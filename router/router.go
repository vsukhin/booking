package router

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vsukhin/booking/controllers"
	"github.com/vsukhin/booking/helpers"
	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/models"
	"github.com/vsukhin/booking/persistence/sqldb"
	"github.com/vsukhin/booking/services"
)

const (
	// APIVersion contains api version
	APIVersion = "v1"
	// version is version regexp
	version = `^\/(v\d\/)?`
	// proxyHeader is proxy header
	proxyHeader = "X-FORWARDED-FOR"
	// ip4Local is ip 4 local
	ip4local = "127.0.0.1"
	// ip6Local is ip 6 local
	ip6local = "::1"
)

// Manager is router manager
type Manager struct {
	db sqldb.DBInterface
}

// ManagerInterface is router manager interface
type ManagerInterface interface {
	InitGin(mode string)
	GinLogger() gin.HandlerFunc
	PanicRecovery() gin.HandlerFunc
	CreateRouter(mode string) *gin.Engine
}

// NewManager is a constructor of router manager
func NewManager(db sqldb.DBInterface) ManagerInterface {
	return &Manager{db: db}
}

func (router *Manager) stackMap(skip int) models.OrderedMap {
	m := models.OrderedMap{}
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		m = append(m, struct {
			Key string
			Val interface{}
		}{
			fmt.Sprintf("%v", i),
			fmt.Sprintf("%s:%d", file, line),
		})
	}

	return m
}

// InitGin inits gin
func (router *Manager) InitGin(mode string) {
	switch mode {
	case logging.ModeDev:
		gin.SetMode(gin.DebugMode)
	case logging.ModeStaging, logging.ModeProd:
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.DebugMode)
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"mode": mode,
		}).Warn("Unexpected mode value")
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"gin": gin.Mode(),
	}).Info("Service gin mode")
}

// GinLogger is gin logger middleware
func (router *Manager) GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		action := ""

		r, err := regexp.Compile(version)
		if err != nil {
			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":  err,
				"path":   path,
				"regexp": version,
			}).Error("Error compiling regexp")
		} else {
			action = r.ReplaceAllString(c.Request.URL.Path, "")
		}

		start := time.Now()
		c.Request.Header.Set("startTime", fmt.Sprintf("%v", start.UnixNano()))

		c.Next()

		stop := time.Since(start)

		var ip = c.Request.Header.Get(proxyHeader)
		if ip == ip6local || ip == ip4local || ip == "" {
			ip, _, err = net.SplitHostPort(c.Request.RemoteAddr)
			if err != nil {
				logging.Log.WithFields(logging.DepthModerate, logging.Fields{
					"error":   err,
					"path":    path,
					"address": c.Request.RemoteAddr,
				}).Error("Error detecting ip")
			}
		}

		m := logging.Fields{
			"action":        action,
			"status_code":   c.Writer.Status(),
			"response_time": fmt.Sprintf("%.4f", stop.Seconds()),
			"ip":            ip,
			"method":        c.Request.Method,
			"request_path":  path,
			"query_string":  c.Request.URL.Query(),
			"user_agent":    c.Request.UserAgent(),
			"created_at":    time.Now().Unix(),
		}

		if len(c.Errors) > 0 {
			m["headers"] = c.Request.Header
			m["stack"] = router.stackMap(logging.DepthLow)
			m["errors"] = c.Errors.String()
			logging.Log.WithFields(logging.DepthModerate, m).Error("request")
		} else {
			logging.Log.WithFields(logging.DepthLow, m).Info("request")
		}
	}
}

// PanicRecovery is panic recovery middleware
func (router *Manager) PanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				m := logging.Fields{
					"stack":      router.stackMap(logging.DepthHigh),
					"error":      err,
					"ip":         c.ClientIP(),
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"user_agent": c.Request.UserAgent(),
				}
				logging.Log.WithFields(logging.DepthModerate, m).Error("Panic recovered")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}

// CreateRouter creates router
func (router *Manager) CreateRouter(mode string) *gin.Engine {
	router.InitGin(mode)

	r := gin.New()

	blockService := services.NewBlockService(router.db)
	seatService := services.NewSeatService(router.db)
	flightService := services.NewFlightService(router.db, blockService, seatService)

	queryManager := helpers.NewQueryManager()

	flightController := controllers.NewFlightController(flightService, queryManager)
	seatController := controllers.NewSeatController(seatService, flightService, queryManager)

	r.Use(router.GinLogger())
	r.Use(router.PanicRecovery())

	v := r.Group("/" + APIVersion)
	{
		v.GET("/flights/:flightId", flightController.Retrieve)
		v.GET("/flights", flightController.ListAll)
		v.OPTIONS("/flights", flightController.GetMeta)
		v.POST("/flights", flightController.Create)
		v.DELETE("/flights/:flightId", flightController.Delete)

		v.GET("/flights/:flightId/seats/index/:index", seatController.Retrieve)
		v.GET("/flights/:flightId/seats/row/:row/line/:line", seatController.Find)
		v.GET("/flights/:flightId/seats", seatController.ListAll)
		v.OPTIONS("/flights/:flightId/seats", seatController.GetMeta)
		v.POST("/flights/:flightId/seats", seatController.Create)
		v.PATCH("/flights/:flightId/seats/:index", seatController.Update)
		v.DELETE("/flights/:flightId/seats/:index", seatController.Delete)
	}

	return r
}
