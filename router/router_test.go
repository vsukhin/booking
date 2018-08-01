package router

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	gorp "gopkg.in/gorp.v2"

	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/persistence/sqldb"
)

func init() {
	logging.Log = NewFakeLogger()
}

// FakeLogger is fake logger
type FakeLogger struct {
	*logrus.Logger
}

// NewFakeLogger is a constructor of fake logger
func NewFakeLogger() logging.LoggerInterface {
	log := logrus.New()

	return &FakeLogger{log}
}

// Init initiates logging
func (logger *FakeLogger) Init(mode string) {
}

// WithFields logs with fields
func (logger *FakeLogger) WithFields(depthLevel int, fields logging.Fields) *logrus.Entry {
	return logrus.NewEntry(logger.Logger)
}

// Info logs info
func (logger *FakeLogger) Info(args ...interface{}) {
}

// FakeDB is fake db management structure
type FakeDB struct {
}

// NewFakeDB creates new fake db management structure
func NewFakeDB() sqldb.DBInterface {
	return &FakeDB{}
}

// AddTableWithName adds table with name to db map
func (db *FakeDB) AddTableWithName(i interface{}, name string) *gorp.TableMap {
	dbmap := &gorp.DbMap{}
	return dbmap.AddTableWithName(i, name)
}

// Insert inserts data to the db table
func (db *FakeDB) Insert(trans *gorp.Transaction, list ...interface{}) error {

	return nil
}

// Update updates data in the db table
func (db *FakeDB) Update(trans *gorp.Transaction, list ...interface{}) (int64, error) {

	return 0, nil
}

// Delete deletes data from the db table
func (db *FakeDB) Delete(trans *gorp.Transaction, list ...interface{}) (int64, error) {

	return 0, nil
}

// Get gets data from the db table
func (db *FakeDB) Get(trans *gorp.Transaction, i interface{}, keys ...interface{}) (interface{}, error) {

	return nil, nil
}

// Select selects data from the db table
func (db *FakeDB) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	var list []interface{}
	return list, nil
}

// SelectInt selects int from the db table
func (db *FakeDB) SelectInt(query string, args ...interface{}) (int64, error) {

	return 0, nil
}

// SelectStr selects string from the db table
func (db *FakeDB) SelectStr(query string, args ...interface{}) (string, error) {

	return "", nil
}

// SelectOne selects one row from the db table
func (db *FakeDB) SelectOne(trans *gorp.Transaction, holder interface{}, query string, args ...interface{}) error {

	return nil
}

// Exec executes statement
func (db *FakeDB) Exec(trans *gorp.Transaction, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

// Begin begins transaction
func (db *FakeDB) Begin() (*gorp.Transaction, error) {
	return nil, nil
}

// Rollback rollbacks transaction
func (db *FakeDB) Rollback(trans *gorp.Transaction) error {
	return nil
}

// Commit commits transaction
func (db *FakeDB) Commit(trans *gorp.Transaction) error {
	return nil
}

// GetDBMap returns dbmap
func (db *FakeDB) GetDBMap() *gorp.DbMap {
	return nil
}

func Test_Router_InitGin_Dev_Success(t *testing.T) {
	router := NewManager(NewFakeDB())

	router.InitGin(logging.ModeDev)
	if gin.Mode() != "debug" {
		t.Error("Expected to have debug mode")
	}
}

func Test_Router_InitGin_Staging_Success(t *testing.T) {
	router := NewManager(NewFakeDB())

	router.InitGin(logging.ModeStaging)
	if gin.Mode() != "release" {
		t.Error("Expected to have release mode")
	}
}

func Test_Router_InitGin_Prod_Success(t *testing.T) {
	router := NewManager(NewFakeDB())

	router.InitGin(logging.ModeProd)
	if gin.Mode() != "release" {
		t.Error("Expected to have release mode")
	}
}

func Test_Router_InitGin_Unknown_Success(t *testing.T) {
	router := NewManager(NewFakeDB())

	router.InitGin("Unknown")
	if gin.Mode() != "debug" {
		t.Error("Expected to have debug mode")
	}
}

func Test_Router_GinLogger_Success(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()

	router := NewManager(NewFakeDB())

	r := gin.New()
	r.Use(router.GinLogger())

	r.POST("/", func(c *gin.Context) {
	})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Expected to return success")
	}
}

func Test_Router_GinLogger_GinError_Failure(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()

	router := NewManager(NewFakeDB())

	r := gin.New()
	r.Use(router.GinLogger())

	r.POST("/", func(c *gin.Context) {
		c.Errors = append(c.Errors, &gin.Error{})
	})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("Expected to return success")
	}
}

func Test_Router_PanicRecovery_Success(t *testing.T) {
	req, _ := http.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()

	router := NewManager(NewFakeDB())

	r := gin.New()
	r.Use(router.PanicRecovery())

	r.POST("/", func(c *gin.Context) {
		panic("test")
	})

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Error("Expected to return server error")
	}
}

func Test_Router_CreateRouter_Success(t *testing.T) {
	router := NewManager(NewFakeDB())

	r := router.CreateRouter(logging.ModeDev)
	if r == nil {
		t.Error("Expected router successfully created")
	}
}
