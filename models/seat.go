package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/vsukhin/booking/logging"
)

// SeatType is seat type
type SeatType int

const (
	// SeatTypeAisle is aisle seat type
	SeatTypeAisle SeatType = iota + 1
	// SeatTypeWindow is window seat type
	SeatTypeWindow
	// SeatTypeMiddle is middle seat type
	SeatTypeMiddle
)

// SeatUpdate is data for seat updating
type SeatUpdate struct {
	Assigned bool `json:"assigned"`
}

// SeatMeta is metadata for seat list
type SeatMeta struct {
	TotalRecords int64 `json:"total_records"`
}

// Seat contains seat data
type Seat struct {
	ID        int64    `json:"id"         db:"id"         query:"id"         search:"id"`
	FlightID  int64    `json:"flight_id"  db:"flight_id"  query:"-"          search:"-"`
	Index     int      `json:"index"      db:"index"      query:"index"      search:"index"`
	Type      SeatType `json:"type"       db:"type"       query:"type"       search:"type"`
	Row       int      `json:"row"        db:"row"        query:"row"        search:"row"`
	Line      string   `json:"line"       db:"line"       query:"line"       search:"line"`
	Assigned  bool     `json:"assigned"   db:"assigned"   query:"assigned"   search:"assigned"`
	CreatedAt int64    `json:"created_at" db:"created_at" query:"created_at" search:"created_at"`
	UpdatedAt int64    `json:"updated_at" db:"updated_at" query:"updated_at" search:"updated_at"`
}

// Validate validates seat data
func (seat *SeatUpdate) Validate() []Error {

	return []Error{}
}

// Verify verifies sort field
func (seat *Seat) Verify(field string) bool {
	return CheckQueryTag(field, seat)
}

// Validate validates search field
func (seat *Seat) Validate(field string, value string) (string, string, []Error) {
	var searchValue string
	var searchField string
	var errs []Error

	searchField = GetSearchTag(field, seat)

	switch field {
	case "id", "created_at", "updated_at":
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
	case "index", "type", "row":
		_, err := strconv.ParseInt(value, 10, 32)
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
	case "line":
		if len(value) != 1 {
			errs = append(errs, Error{
				Code:    field + ".Invalid",
				Message: field + " is not one character",
				Field:   field,
			})

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"errors": errs,
				"field":  field,
				"value":  value,
			}).Error(field + " is not one character")
			break
		}

		if strings.Contains(value, "'") {
			value = strings.Replace(value, "'", "''", -1)
		}
		searchValue = "'" + value + "'"
	case "assigned":
		val, err := strconv.ParseBool(value)
		if err != nil {
			errs = append(errs, Error{
				Code:    field + ".Invalid",
				Message: field + " is not boolean",
				Field:   field,
			})

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":  err,
				"errors": errs,
				"field":  field,
				"value":  value,
			}).Error(field + " is not boolean")
			break
		}
		searchValue = fmt.Sprintf("%v", val)
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
func (seat *Seat) GetAllFields() []string {
	return GetAllSearchTags(seat)
}

// ValidateAll validates all seach fields
func (seat *Seat) ValidateAll(value string) string {
	var searchValue string

	if strings.Contains(value, "'") {
		value = strings.Replace(value, "'", "''", -1)
	}

	searchValue = "'" + value + "'"
	return searchValue
}
