package models

const (
	// maxFieldLength is max field length
	maxFieldLength = 255
)

// Error contains error details
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field"`
}
