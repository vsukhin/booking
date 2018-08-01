package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vsukhin/booking/logging"
)

const (
	// maxLines is max line number
	maxLines = 20
	// maxRows is max rows number
	maxRows = 200
)

// FlightCreate is data for flight creation
type FlightCreate struct {
	Name   string  `json:"name"`
	Blocks []Block `json:"blocks"`
}

// FlightMeta is metadata for flight list
type FlightMeta struct {
	TotalRecords int64 `json:"total_records"`
}

// Flight contains flight data
type Flight struct {
	ID        int64   `json:"id"               db:"id"         query:"id"         search:"id"`
	Name      string  `json:"name"             db:"name"       query:"name"       search:"name"`
	CreatedAt int64   `json:"created_at"       db:"created_at" query:"created_at" search:"created_at"`
	Blocks    []Block `json:"blocks,omitempty" db:"-"`
}

// Validate validates flight data
func (flight *FlightCreate) Validate() []Error {
	var errs []Error

	if len([]rune(flight.Name)) > maxFieldLength {
		errs = append(errs, Error{
			Code:    "name.TooLarge",
			Message: fmt.Sprintf("Name must be less than %v characters", maxFieldLength),
			Field:   "name",
		})
	}

	rows := 0
	for _, block := range flight.Blocks {
		valerrs := block.Validate()
		errs = append(errs, valerrs...)
		rows += block.Rows
	}

	if rows > maxRows {
		errs = append(errs, Error{
			Code:    "rows.TooLarge",
			Message: fmt.Sprintf("Rows must be less than %v", maxRows),
			Field:   "rows",
		})
	}

	return errs
}

// Verify verifies sort field
func (flight *Flight) Verify(field string) bool {
	return CheckQueryTag(field, flight)
}

// Validate validates search field
func (flight *Flight) Validate(field string, value string) (string, string, []Error) {
	var searchValue string
	var searchField string
	var errs []Error

	searchField = GetSearchTag(field, flight)

	switch field {
	case "id", "created_at":
		_, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			errs = append(errs, Error{
				Code:    field + ".Invalid",
				Message: field + " is not integer",
				Field:   field,
			})

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":  err,
				"errors": errs,
				"field":  field,
				"value":  value,
			}).Error(field + " is not integer")
			break
		}
		searchValue = value
	case "name":
		if strings.Contains(value, "'") {
			value = strings.Replace(value, "'", "''", -1)
		}
		searchValue = "'" + value + "'"
	default:
		errs = append(errs, Error{
			Code:    "field.Unknown",
			Message: field + " field is unknown",
			Field:   "field",
		})

		logging.Log.WithFields(logging.DepthModerate, logging.Fields{
			"errors": errs,
			"field":  field,
			"value":  value,
		}).Error(field + " field is unknown")
	}

	return searchField, searchValue, errs
}

// GetAllFields gets all search fields
func (flight *Flight) GetAllFields() []string {
	return GetAllSearchTags(flight)
}

// ValidateAll validates all seach fields
func (flight *Flight) ValidateAll(value string) string {
	var searchValue string

	if strings.Contains(value, "'") {
		value = strings.Replace(value, "'", "''", -1)
	}

	searchValue = "'" + value + "'"
	return searchValue
}
