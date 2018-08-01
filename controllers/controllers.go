package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/models"
	"github.com/vsukhin/booking/services"
)

func getFlight(c *gin.Context, flightService services.FlightServiceInterface) (*models.Flight, error) {
	flightID, err := strconv.ParseInt(c.Params.ByName("flightId"), 10, 64)
	if err != nil {
		errs := []models.Error{models.Error{
			Code:    "flightId.Invalid",
			Message: "Flight id is not integer",
			Field:   "flightId",
		}}

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"error":    err,
			"errors":   errs,
			"flightId": c.Params.ByName("flightId"),
		}).Error("Flight id is not integer")

		c.JSON(http.StatusBadRequest, errs)
		return nil, err
	}

	flight, err := flightService.Retrieve(flightID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return nil, err
	}

	if flight == nil {
		c.Status(http.StatusNotFound)
		return nil, errors.New("Flight not found")
	}

	return flight, nil
}
