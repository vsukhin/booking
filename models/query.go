package models

// SortFieldChecker is sort field checker interface
type SortFieldChecker interface {
	Verify(field string) bool
}

// OrderExp is order expression
type OrderExp struct {
	Field string
	Order string
}

// SearchFieldChecker is search field checker interface
type SearchFieldChecker interface {
	Validate(field string, value string) (string, string, []Error)
	ValidateAll(value string) string
	GetAllFields() []string
}

// FilterExp is filter expression
type FilterExp struct {
	Fields []string
	Op     string
	Value  string
}
