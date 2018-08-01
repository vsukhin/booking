package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vsukhin/booking/helpers"
	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/models"
	"github.com/vsukhin/booking/services"
)

// SeatController is an seat controller
type SeatController struct {
	seatService   services.SeatServiceInterface
	flightService services.FlightServiceInterface
	queryManager  helpers.QueryManagerInterface
}

// SeatControllerInterface is an interface for seat controller methods
type SeatControllerInterface interface {
	Retrieve(c *gin.Context)
	Find(c *gin.Context)
	ListAll(c *gin.Context)
	GetMeta(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

// NewSeatController is a constructor for seat controller
func NewSeatController(seatService services.SeatServiceInterface, flightService services.FlightServiceInterface,
	queryManager helpers.QueryManagerInterface) SeatControllerInterface {
	return &SeatController{seatService: seatService, flightService: flightService, queryManager: queryManager}
}

func (seatController *SeatController) getSeat(c *gin.Context) (*models.Seat, error) {
	flight, err := getFlight(c, seatController.flightService)
	if err != nil {
		return nil, err
	}

	index, err := strconv.ParseInt(c.Params.ByName("index"), 10, 32)
	if err != nil {
		errs := []models.Error{models.Error{
			Code:    "index.Invalid",
			Message: "Index is not integer",
			Field:   "index",
		}}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":  err,
			"errors": errs,
			"seatId": c.Params.ByName("index"),
		}).Error("Index is not integer")

		c.JSON(http.StatusBadRequest, errs)
		return nil, err
	}

	seat, err := seatController.seatService.Retrieve(flight.ID, index)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return nil, err
	}

	if seat == nil {
		c.Status(http.StatusNotFound)
		return nil, errors.New("Seat not found")
	}

	return seat, nil
}

// Retrieve retrieves seat
func (seatController *SeatController) Retrieve(c *gin.Context) {
	seat, err := seatController.getSeat(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, seat)
}

// Find finds seat
func (seatController *SeatController) Find(c *gin.Context) {
	flight, err := getFlight(c, seatController.flightService)
	if err != nil {
		return
	}

	row, err := strconv.Atoi(c.Params.ByName("row"))
	if err != nil {
		errs := []models.Error{models.Error{
			Code:    "row.Invalid",
			Message: "Row is not integer",
			Field:   "row",
		}}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":  err,
			"errors": errs,
			"row":    c.Params.ByName("row"),
		}).Error("Row is not integer")

		c.JSON(http.StatusBadRequest, errs)
		return
	}

	line := c.Params.ByName("line")
	if len(line) != 1 {
		errs := []models.Error{models.Error{
			Code:    "line.Invalid",
			Message: "Line is not one character",
			Field:   "line",
		}}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"errors": errs,
			"line":   c.Params.ByName("line"),
		}).Error("Line is not one character")

		c.JSON(http.StatusBadRequest, errs)
		return
	}

	seat, err := seatController.seatService.Find(flight.ID, row, line)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if seat == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, seat)
}

// ListAll lists all seats according filter, sort, offset, limit parameters
func (seatController *SeatController) ListAll(c *gin.Context) {
	flight, err := getFlight(c, seatController.flightService)
	if err != nil {
		return
	}

	limitation, errs := seatController.queryManager.GetLimitation(c)
	if len(errs) != 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	sorting, errs := seatController.queryManager.GetSorting(&models.Seat{}, c)
	if len(errs) != 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	filtering, errs := seatController.queryManager.GetFiltering(&models.Seat{}, c)
	if len(errs) != 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	seats, err := seatController.seatService.ListAll(flight.ID, filtering, sorting, limitation)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, seats)
}

// GetMeta gets meta data about seat list according filter parameters
func (seatController *SeatController) GetMeta(c *gin.Context) {
	flight, err := getFlight(c, seatController.flightService)
	if err != nil {
		return
	}

	filtering, errs := seatController.queryManager.GetFiltering(&models.Seat{}, c)
	if len(errs) != 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	seatMeta, err := seatController.seatService.GetMeta(flight.ID, filtering)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, seatMeta)
}

// Create creates seat
func (seatController *SeatController) Create(c *gin.Context) {
	flight, err := getFlight(c, seatController.flightService)
	if err != nil {
		return
	}

	seat, err := seatController.seatService.Assign(flight.ID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, seat)
}

// Update updates seat
func (seatController *SeatController) Update(c *gin.Context) {
	seat, err := seatController.getSeat(c)
	if err != nil {
		return
	}

	var seatUpdate models.SeatUpdate

	err = c.BindJSON(&seatUpdate)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error": err,
		}).Error("Error binding seat")

		c.Status(http.StatusBadRequest)
		return
	}

	errs := seatUpdate.Validate()
	if len(errs) != 0 {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"seatUpdate": seatUpdate,
			"errors":     errs,
		}).Error("Error validating seat")

		c.JSON(http.StatusBadRequest, errs)
		return
	}

	seat.Assigned = seatUpdate.Assigned
	seat.UpdatedAt = time.Now().Unix()

	err = seatController.seatService.Update(seat)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, seat)
}

// Delete deletes seat
func (seatController *SeatController) Delete(c *gin.Context) {
	seat, err := seatController.getSeat(c)
	if err != nil {
		return
	}

	seat.Assigned = false
	seat.UpdatedAt = time.Now().Unix()

	err = seatController.seatService.Update(seat)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
