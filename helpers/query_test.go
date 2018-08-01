package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/vsukhin/booking/logging"
	"github.com/vsukhin/booking/models"
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

func Test_GetLimitation_Success(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?offset=5&limit=10", nil)

	manager := NewQueryManager()

	limitation, errs := manager.GetLimitation(c)
	if len(errs) != 0 {
		t.Error("Expected to get limitation successfully")
	}
	if limitation != " LIMIT 5, 10" {
		t.Error("Expected to get matching limitation")
	}
}

func Test_GetLimitation_Empty_Success(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)

	manager := NewQueryManager()

	limitation, errs := manager.GetLimitation(c)
	if len(errs) != 0 {
		t.Error("Expected to get empty limitation successfully")
	}
	if limitation != " LIMIT 0, 100" {
		t.Error("Expected to get matching limitation")
	}
}

func Test_GetLimitation_OffsetInteger_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?offset=a&limit=1", nil)

	manager := NewQueryManager()

	limitation, errs := manager.GetLimitation(c)
	if len(errs) != 1 {
		t.Error("Expected to have getting limitation not integer error")
	}
	for _, err := range errs {
		if err.Code != "offset.Invalid" {
			t.Error("Expected to have getting limitation not integer offset error")
		}
	}
	if limitation != "" {
		t.Error("Expected to have empty limitation")
	}
}

func Test_GetLimitation_LimitInteger_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?offset=1&limit=a", nil)

	manager := NewQueryManager()

	limitation, errs := manager.GetLimitation(c)
	if len(errs) != 1 {
		t.Error("Expected to have getting limitation not integer error")
	}
	for _, err := range errs {
		if err.Code != "limit.Invalid" {
			t.Error("Expected to have getting limitation not integer limit error")
		}
	}
	if limitation != "" {
		t.Error("Expected to have empty limitation")
	}
}

func Test_GetLimitation_OffsetNegative_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?offset=-1&limit=1", nil)

	manager := NewQueryManager()

	limitation, errs := manager.GetLimitation(c)
	if len(errs) != 1 {
		t.Error("Expected to have getting limitation negative error")
	}
	for _, err := range errs {
		if err.Code != "offset.Negative" {
			t.Error("Expected to have getting limitation negative offset error")
		}
	}
	if limitation != "" {
		t.Error("Expected to have empty limitation")
	}
}

func Test_GetLimitation_LimitNegative_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?offset=1&limit=-1", nil)

	manager := NewQueryManager()

	limitation, errs := manager.GetLimitation(c)
	if len(errs) != 1 {
		t.Error("Expected to have getting limitation negative error")
	}
	for _, err := range errs {
		if err.Code != "limit.Negative" {
			t.Error("Expected to have getting limitation negative limit error")
		}
	}
	if limitation != "" {
		t.Error("Expected to have empty limitation")
	}
}

func Test_GetSorting_Success(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?sort=line:asc&sort=id:desc&sort=", nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	sorting, errs := manager.GetSorting(seat, c)
	if len(errs) != 0 {
		t.Error("Expected to get sorting successfully")
	}
	if sorting != " ORDER BY line ASC, id DESC" {
		t.Error("Expected to get matching sorting")
	}
}

func Test_GetSorting_Split_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?sort=line", nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	sorting, errs := manager.GetSorting(seat, c)
	if len(errs) != 1 {
		t.Error("Expected to have getting sorting splitting error")
	}
	for _, err := range errs {
		if err.Code != "sort.WrongLength" {
			t.Error("Expected to have getting sorting splitting length error")
		}
	}
	if sorting != "" {
		t.Error("Expected to have empty sorting")
	}
}

func Test_GetSorting_Field_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?sort=test:asc", nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	sorting, errs := manager.GetSorting(seat, c)
	if len(errs) != 1 {
		t.Error("Expected to have getting sorting field error")
	}
	for _, err := range errs {
		if err.Code != "sort.UnknownField" {
			t.Error("Expected to have getting sorting field name error")
		}
	}
	if sorting != "" {
		t.Error("Expected to have empty sorting")
	}
}

func Test_GetSorting_Order_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?sort=id:test", nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	sorting, errs := manager.GetSorting(seat, c)
	if len(errs) != 1 {
		t.Error("Expected to have getting sorting order error")
	}
	for _, err := range errs {
		if err.Code != "sort.UnknownOrder" {
			t.Error("Expected to have getting sorting order value error")
		}
	}
	if sorting != "" {
		t.Error("Expected to have empty sorting")
	}
}

func Test_GetFiltering_Success(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET",
		"/?filter=line:eq:v&filter=id:gt:100&filter=id:lt:2&filter=id:le:3&filter=id:ge:4&filter=id:ne:5&filter=id:lk:6",
		nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	filtering, errs := manager.GetFiltering(seat, c)
	if len(errs) != 0 {
		t.Error("Expected to get filtering successfully")
	}
	if filtering != " AND (line = 'v') AND (id > 100) AND (id < 2) AND (id <= 3) AND (id >= 4) AND (id != 5) AND (id LIKE 6)" {
		t.Error("Expected to get matching filtering")
	}
}

func Test_GetFiltering_Empty_Success(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?filter=", nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	filtering, errs := manager.GetFiltering(seat, c)
	if len(errs) != 0 {
		t.Error("Expected to get filtering successfully")
	}
	if filtering != "" {
		t.Error("Expected to get empty filtering")
	}
}

func Test_GetFiltering_Spilt_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?filter=line:eq", nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	filtering, errs := manager.GetFiltering(seat, c)
	if len(errs) != 1 {
		t.Error("Expected to have getting filtering splitting error")
	}
	for _, err := range errs {
		if err.Code != "filter.WrongLength" {
			t.Error("Expected to have getting filtering splitting length error")
		}
	}
	if filtering != "" {
		t.Error("Expected to have empty filtering")
	}
}

func Test_GetFiltering_Operation_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?filter=line:test:1", nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	filtering, errs := manager.GetFiltering(seat, c)
	if len(errs) != 1 {
		t.Error("Expected to have getting filtering operation error")
	}
	for _, err := range errs {
		if err.Code != "filter.UnknownOperation" {
			t.Error("Expected to have getting filtering operation name error")
		}
	}
	if filtering != "" {
		t.Error("Expected to have empty filtering")
	}
}

func Test_GetFiltering_Value_Failure(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?filter=id:eq:test", nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	filtering, errs := manager.GetFiltering(seat, c)
	if len(errs) != 1 {
		t.Error("Expected to have getting filtering value error")
	}
	for _, err := range errs {
		if err.Code != "id.Invalid" {
			t.Error("Expected to have getting filtering value not integer error")
		}
	}
	if filtering != "" {
		t.Error("Expected to have empty filtering")
	}
}

func Test_GetFiltering_All_Success(t *testing.T) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/?filter=*:lk:\"*v' s\"", nil)

	manager := NewQueryManager()
	seat := &models.Seat{}

	filtering, errs := manager.GetFiltering(seat, c)
	if len(errs) != 0 {
		t.Error("Expected to get all filtering successfully")
	}
	if filtering != " AND (id LIKE '%v'' s' OR index LIKE '%v'' s' OR type LIKE '%v'' s' "+
		"OR row LIKE '%v'' s' OR line LIKE '%v'' s' OR assigned LIKE '%v'' s' OR created_at LIKE '%v'' s' "+
		"OR updated_at LIKE '%v'' s')" {
		t.Error("Expected to get matching filtering")
	}
}
