package handlers

import (
	"bytes"
	"fizzbuzz-api/internal/fizzbuzzapi/config"
	"fizzbuzz-api/internal/fizzbuzzapi/types"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/* Mock implementations for testing */

type mockController struct {
	types.FizzBuzzLimits
}

func (m *mockController) GenerateFizzBuzz(req types.FizzBuzzRequest) (types.FizzBuzzResponse, error) {
	return types.FizzBuzzResponse{}, nil
}

type mockRecorder struct{}

func (m *mockRecorder) GetStats() types.FizzBuzzStats {
	return types.FizzBuzzStats{}
}
func (m *mockRecorder) SaveStat(req types.FizzBuzzRequest) error {
	return nil
}

type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...any)  {}
func (m *mockLogger) Error(msg string, args ...any) {}
func (m *mockLogger) Debug(msg string, args ...any) {}

var mockFizzBuzzController = &mockController{}
var mockStatsRecorder = &mockRecorder{}
var mockConfig = config.Config{
	MaxFizzBuzzLimit: 100,
	MaxStringLength:  100,
}

/* Test functions */

func initMockGinRequest(body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	req := httptest.NewRequest("POST", "/fizzbuzz/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	return c, w
}

func Test_GenerateFizzBuzz(t *testing.T) {
	assert := assert.New(t)

	body := []byte(`{"int1":1,"int2":2,"limit":2,"str1":"a","str2":"b"}`)
	c, w := initMockGinRequest(body)

	handler := NewFizzBuzzHandler(&mockConfig, &mockLogger{}, mockFizzBuzzController, mockStatsRecorder)
	assert.NotNil(handler)

	handler.GenerateFizzBuzz(c)
	assert.Equal(200, w.Code)
}

func Test_GenerateFizzBuzz_InvalidRequestTypes(t *testing.T) {
	assert := assert.New(t)
	body := []byte(`{"int1":"invalid","int2":2,"limit":2,"str1":"a","str2":"b"}`)
	c, w := initMockGinRequest(body)
	handler := NewFizzBuzzHandler(&mockConfig, &mockLogger{}, mockFizzBuzzController, mockStatsRecorder)
	assert.NotNil(handler)
	handler.GenerateFizzBuzz(c)
	assert.Equal(400, w.Code)
}

func Test_GenerateFizzBuzz_MissingRequestParams(t *testing.T) {
	assert := assert.New(t)
	body := []byte(`{"int2":2,"limit":2,"str1":"a","str2":"b"}`)
	c, w := initMockGinRequest(body)
	handler := NewFizzBuzzHandler(&mockConfig, &mockLogger{}, mockFizzBuzzController, mockStatsRecorder)
	assert.NotNil(handler)
	handler.GenerateFizzBuzz(c)
	assert.Equal(400, w.Code)
}
