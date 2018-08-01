package helpers

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/models"
)

const (
	// queryParameterOffset is offset query parameter
	queryParameterOffset = "offset"
	// queryParameterLimit is limit query parameter
	queryParameterLimit = "limit"
	// queryParameterSort is sort query parameter
	queryParameterSort = "sort"
	// queryParameterFilter is filter query parameter
	queryParameterFilter = "filter"

	// indexOffset is offset index
	indexOffset = 0
	// indexLimit is limit index
	indexLimit = 1
	// defaultLimit is default limit
	defaultLimit = 100

	// queryParameterSortAsc is ascending sort query parameter
	queryParameterSortAsc = "asc"
	// queryParameterSortDesc is descending sort query parameter
	queryParameterSortDesc = "desc"
	// queryParameterSortField is sort field query parameter
	queryParameterSortField = 0
	// queryParameterSortOrder is sort order query parameter
	queryParameterSortOrder = 1
	// queryParameterSortLength is length of query sort parameters
	queryParameterSortLength = 2

	// queryParameterFilterOpEq is filter == operation parameter
	queryParameterFilterOpEq = "eq"
	// queryParameterFilterOpNe is filter != operation parameter
	queryParameterFilterOpNe = "ne"
	// queryParameterFilterOpLt is filter < operation parameter
	queryParameterFilterOpLt = "lt"
	// queryParameterFilterOpLe is filter <= operation parameter
	queryParameterFilterOpLe = "le"
	// queryParameterFilterOpGt is filter > operation parameter
	queryParameterFilterOpGt = "gt"
	// queryParameterFilterOpGe is filter >= operation parameter
	queryParameterFilterOpGe = "ge"
	// queryParameterFilterOpLk is filter like operation parameter
	queryParameterFilterOpLk = "lk"
	// queryParameterFilterField is filter field query parameter
	queryParameterFilterField = 0
	// queryParameterFilterOp is filter operand query parameter
	queryParameterFilterOp = 1
	// queryParameterFilterValue is filter value query parameter
	queryParameterFilterValue = 2
	// queryParameterFilterLength is length of query filter parameters
	queryParameterFilterLength = 3

	// delimiter is expression operand delimiter
	delimiter = ':'
)

// QueryManager is query manager
type QueryManager struct {
}

// QueryManagerInterface is query manager interface
type QueryManagerInterface interface {
	GetLimitation(c *gin.Context) (string, []models.Error)
	GetSorting(checker models.SortFieldChecker, c *gin.Context) (string, []models.Error)
	GetFiltering(checker models.SearchFieldChecker, c *gin.Context) (string, []models.Error)
}

// NewQueryManager is a constructor of query manager
func NewQueryManager() QueryManagerInterface {
	return &QueryManager{}
}

// GetLimitation get limitation from the query
func (manager *QueryManager) GetLimitation(c *gin.Context) (string, []models.Error) {
	var limitation string
	var fields = []struct {
		data    string
		name    string
		message string
		value   int64
	}{
		{
			queryParameterOffset,
			"offset",
			"Offset",
			0,
		},
		{
			queryParameterLimit,
			"limit",
			"Limit",
			0,
		},
	}

	for i := range fields {
		value, err := url.QueryUnescape(c.Request.URL.Query().Get(fields[i].data))
		if err != nil {
			errs := []models.Error{models.Error{
				Code:    fields[i].name + ".Bad",
				Message: fields[i].message + " can't be unescaped",
				Field:   fields[i].name,
			}}

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":  err,
				"errors": errs,
				"query":  c.Request.URL.RawQuery,
			}).Error(fields[i].message + " can't be unescaped")

			return "", errs
		}

		if value != "" {
			var valueInt int64

			valueInt, err = strconv.ParseInt(value, 10, 64)
			if err != nil {
				errs := []models.Error{models.Error{
					Code:    fields[i].name + ".Invalid",
					Message: fields[i].message + " is not integer",
					Field:   fields[i].name,
				}}

				logging.Log.WithFields(logging.DepthModerate, logging.Fields{
					"error":        err,
					"errors":       errs,
					fields[i].name: value,
					"query":        c.Request.URL.RawQuery,
				}).Error(fields[i].message + " is not integer")

				return "", errs
			}

			if valueInt < 0 {
				errs := []models.Error{models.Error{
					Code:    fields[i].name + ".Negative",
					Message: fields[i].message + " can't be negative",
					Field:   fields[i].name,
				}}

				logging.Log.WithFields(logging.DepthModerate, logging.Fields{
					"errors":       errs,
					fields[i].name: value,
					"query":        c.Request.URL.RawQuery,
				}).Error(fields[i].message + " can't be negative")

				return "", errs
			}

			fields[i].value = valueInt
		}
	}

	if fields[indexLimit].value > 0 {
		limitation = fmt.Sprintf(" LIMIT %v, %v", fields[indexOffset].value, fields[indexLimit].value)
	} else {
		limitation = fmt.Sprintf(" LIMIT %v, %v", fields[indexOffset].value, defaultLimit)
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"limitation": limitation,
		"query":      c.Request.URL.RawQuery,
	}).Debug("Offset and limit successfully parsed")
	return limitation, []models.Error{}
}

// GetSorting gets sorting from the query
func (manager *QueryManager) GetSorting(checker models.SortFieldChecker, c *gin.Context) (string, []models.Error) {
	var sorting string

	orders := c.Request.URL.Query()[queryParameterSort]

	var sorts []models.OrderExp

	for _, order := range orders {
		element, err := url.QueryUnescape(order)
		if err != nil {
			errs := []models.Error{models.Error{
				Code:    "sort.Bad",
				Message: "Sort can't be unescaped",
				Field:   "sort",
			}}

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":  err,
				"errors": errs,
				"order":  order,
				"query":  c.Request.URL.RawQuery,
			}).Error("Sort can't be unescaped")

			return "", errs
		}

		if element == "" {
			continue
		}

		elements := strings.Split(element, string(delimiter))
		if len(elements) != queryParameterSortLength {
			errs := []models.Error{models.Error{
				Code:    "sort.WrongLength",
				Message: "Sort has wrong length of elements",
				Field:   "sort",
			}}

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"errors":  errs,
				"element": element,
				"query":   c.Request.URL.RawQuery,
			}).Error("Sort has wrong length of elements")

			return "", errs
		}

		fieldElelemnt := elements[queryParameterSortField]
		orderElement := elements[queryParameterSortOrder]

		valid := checker.Verify(fieldElelemnt)
		if !valid {
			errs := []models.Error{models.Error{
				Code:    "sort.UnknownField",
				Message: "Sort contains unknown field",
				Field:   "sort",
			}}

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"errors":       errs,
				"element":      element,
				"fieldElement": fieldElelemnt,
				"query":        c.Request.URL.RawQuery,
			}).Error("Sort contains unknown field")

			return "", errs
		}

		if strings.ToLower(orderElement) != queryParameterSortAsc &&
			strings.ToLower(orderElement) != queryParameterSortDesc {
			errs := []models.Error{models.Error{
				Code:    "sort.UnknownOrder",
				Message: "Sort contains unknown order",
				Field:   "sort",
			}}

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"errors":       errs,
				"element":      element,
				"orderElement": orderElement,
				"query":        c.Request.URL.RawQuery,
			}).Error("Sort contains unknown order")

			return "", errs
		}

		sorts = append(sorts, models.OrderExp{
			Field: fieldElelemnt,
			Order: strings.ToUpper(orderElement),
		})
	}

	if len(sorts) != 0 {
		var orders []string

		for _, sort := range sorts {
			orders = append(orders, " "+sort.Field+" "+sort.Order)
		}

		sorting += " ORDER BY"
		sorting += strings.Join(orders, ",")
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"sorting": sorting,
		"query":   c.Request.URL.RawQuery,
	}).Debug("Sort successfully parsed")
	return sorting, []models.Error{}
}

// GetFiltering gets filtering from the query
func (manager *QueryManager) GetFiltering(checker models.SearchFieldChecker, c *gin.Context) (string, []models.Error) {
	var filtering string

	expressions := c.Request.URL.Query()[queryParameterFilter]

	var filters []models.FilterExp

	for _, expression := range expressions {
		element, err := url.QueryUnescape(expression)
		if err != nil {
			errs := []models.Error{models.Error{
				Code:    "filter.Bad",
				Message: "Filter can't be unescaped",
				Field:   "filter",
			}}

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":      err,
				"errors":     errs,
				"expression": expression,
				"query":      c.Request.URL.RawQuery,
			}).Error("Filter can't be unescaped")

			return "", errs
		}

		r := csv.NewReader(strings.NewReader(element))
		r.Comma = delimiter
		r.LazyQuotes = true

		elements, err := r.ReadAll()
		if err != nil {
			errs := []models.Error{models.Error{
				Code:    "filter.InvalidFormat",
				Message: "Filter is not in csv format",
				Field:   "filter",
			}}

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"error":   err,
				"errors":  errs,
				"element": element,
				"query":   c.Request.URL.RawQuery,
			}).Error("Filter is not in csv format")

			return "", errs
		}

		if len(elements) == 0 {
			continue
		}

		if len(elements[0]) != queryParameterFilterLength {
			errs := []models.Error{models.Error{
				Code:    "filter.WrongLength",
				Message: "Filter has wrong length of elements",
				Field:   "filter",
			}}

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"errors":  errs,
				"element": element,
				"query":   c.Request.URL.RawQuery,
			}).Error("Filter has wrong length of elements")

			return "", errs
		}

		var allFields bool
		var field string
		var value string

		fieldElement := elements[0][queryParameterFilterField]
		opElement := elements[0][queryParameterFilterOp]
		valueElement := elements[0][queryParameterFilterValue]

		if fieldElement == "*" {
			allFields = true
		}

		if allFields {
			field = ""
			value = checker.ValidateAll(valueElement)
		} else {
			var errs []models.Error

			field, value, errs = checker.Validate(fieldElement, valueElement)
			if len(errs) != 0 {
				return "", errs
			}
		}

		op := ""
		switch strings.ToLower(opElement) {
		case queryParameterFilterOpEq:
			op = "="
		case queryParameterFilterOpLt:
			op = "<"
		case queryParameterFilterOpLe:
			op = "<="
		case queryParameterFilterOpGt:
			op = ">"
		case queryParameterFilterOpGe:
			op = ">="
		case queryParameterFilterOpNe:
			op = "!="
		case queryParameterFilterOpLk:
			op = "LIKE"
			value = strings.Replace(value, "*", "%", -1)
		default:
			errs := []models.Error{models.Error{
				Code:    "filter.UnknownOperation",
				Message: "Filter contains unknown operation",
				Field:   "filter",
			}}

			logging.Log.WithFields(logging.DepthModerate, logging.Fields{
				"errors":    errs,
				"element":   element,
				"opElement": opElement,
				"query":     c.Request.URL.RawQuery,
			}).Error("Filter contains unknown operation")

			return "", errs
		}

		if allFields {
			var fields []string

			for _, field = range checker.GetAllFields() {
				fields = append(fields, field)
			}

			filters = append(filters, models.FilterExp{
				Fields: fields,
				Op:     op,
				Value:  value,
			})
		} else {
			filters = append(filters, models.FilterExp{
				Fields: []string{field},
				Op:     op,
				Value:  value,
			})
		}
	}

	if len(filters) != 0 {
		var masks []string

		for _, filter := range filters {
			var exps []string

			for _, field := range filter.Fields {
				exps = append(exps, field+" "+filter.Op+" "+filter.Value)
			}

			masks = append(masks, "("+strings.Join(exps, " OR ")+")")
		}

		filtering += " AND "
		filtering += strings.Join(masks, " AND ")
	}

	logging.Log.WithFields(logging.DepthLow, logging.Fields{
		"filtering": filtering,
		"query":     c.Request.URL.RawQuery,
	}).Debug("Filter successfully parsed")
	return filtering, []models.Error{}
}
