package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vsukhin/booking/helpers"
	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/models"
	"github.com/vsukhin/booking/services"
)

// FlightController is an flight controller
type FlightController struct {
	flightService services.FlightServiceInterface
	queryManager  helpers.QueryManagerInterface
}

// FlightControllerInterface is an interface for flight controller methods
type FlightControllerInterface interface {
	Retrieve(c *gin.Context)
	ListAll(c *gin.Context)
	GetMeta(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
}

// NewFlightController is a constructor for flight controller
func NewFlightController(flightService services.FlightServiceInterface,
	queryManager helpers.QueryManagerInterface) FlightControllerInterface {
	return &FlightController{flightService: flightService, queryManager: queryManager}
}

// Retrieve retrieves flight
func (flightController *FlightController) Retrieve(c *gin.Context) {
	flight, err := getFlight(c, flightController.flightService)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, flight)
}

// ListAll lists all flights according filter, sort, offset, limit parameters
func (flightController *FlightController) ListAll(c *gin.Context) {
	limitation, errs := flightController.queryManager.GetLimitation(c)
	if len(errs) != 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	sorting, errs := flightController.queryManager.GetSorting(&models.Flight{}, c)
	if len(errs) != 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	filtering, errs := flightController.queryManager.GetFiltering(&models.Flight{}, c)
	if len(errs) != 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	flights, err := flightController.flightService.ListAll(filtering, sorting, limitation)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, flights)
}

// GetMeta gets meta data about flight list according filter parameters
func (flightController *FlightController) GetMeta(c *gin.Context) {
	filtering, errs := flightController.queryManager.GetFiltering(&models.Flight{}, c)
	if len(errs) != 0 {
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	flightMeta, err := flightController.flightService.GetMeta(filtering)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, flightMeta)
}

// Create creates flight
func (flightController *FlightController) Create(c *gin.Context) {
	var flightCreate models.FlightCreate

	err := c.BindJSON(&flightCreate)
	if err != nil {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error": err,
		}).Error("Error binding flight")

		c.Status(http.StatusBadRequest)
		return
	}

	errs := flightCreate.Validate()
	if len(errs) != 0 {
		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"flightCreate": flightCreate,
			"errors":       errs,
		}).Error("Error validating flight")

		c.JSON(http.StatusBadRequest, errs)
		return
	}

	flight := &models.Flight{
		Name:      flightCreate.Name,
		Blocks:    flightCreate.Blocks,
		CreatedAt: time.Now().Unix(),
	}

	err = flightController.flightService.Create(flight)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, flight)
}

// Delete deletes flight
func (flightController *FlightController) Delete(c *gin.Context) {
	flight, err := getFlight(c, flightController.flightService)
	if err != nil {
		return
	}

	err = flightController.flightService.Delete(flight)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
