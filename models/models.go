package models

import (
	"reflect"
)

const (
	// QueryTag is query tag
	QueryTag = "query"
	// SearchTag is search tag
	SearchTag = "search"
	// DateFormat is date format
	DateFormat = "01/02/2006"
	// DateTimeFormat is date time format
	DateTimeFormat = "01/02/2006 15:04"
)

// CheckQueryTag checks availability of query tag
func CheckQueryTag(field string, object interface{}) bool {
	var found bool

	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get(QueryTag)
		if field == fieldTag {
			found = true
			break
		}
	}

	return found
}

// GetSearchTag gets search tag
func GetSearchTag(field string, object interface{}) string {
	var search string

	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get(QueryTag)
		if field == fieldTag {
			search = structAddr.Type().Field(i).Tag.Get(SearchTag)
			break
		}
	}

	return search
}

// GetAllSearchTags gets all search tags
func GetAllSearchTags(object interface{}) []string {
	var tags []string

	structAddr := reflect.ValueOf(object).Elem()
	for i := 0; i < structAddr.NumField(); i++ {
		fieldTag := structAddr.Type().Field(i).Tag.Get(SearchTag)
		if fieldTag != "" && fieldTag != "-" {
			tags = append(tags, fieldTag)
		}
	}

	return tags
}
