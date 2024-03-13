package api

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sitetester/sequence-api/api"
	"github.com/sitetester/sequence-api/config"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func checkStatusCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}

func parseErrorResponse(r *httptest.ResponseRecorder) *api.ErrorResponse {
	var er api.ErrorResponse
	json.NewDecoder(r.Body).Decode(&er)
	return &er
}

var Db *gorm.DB = nil
var engine *gin.Engine = nil

// let's setup DB & router once
func setupTestEnv() {
	gin.SetMode(gin.TestMode) // switch to test mode (to avoid debug output)

	if Db == nil {
		Db = config.SetupDb("../../db/sequences_test.db")
	}

	if engine == nil {
		engine = config.SetupRouter(Db)
	}
}

func performRequest(t *testing.T, method string, url string, data any) *httptest.ResponseRecorder {
	body, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Couldn't marshal JSON: %v\n", err)
	}

	// mock request
	request, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// create a response recorder so can inspect the response
	recorder := httptest.NewRecorder()

	// perform the request
	engine.ServeHTTP(recorder, request)
	return recorder
}

func checkFailsWithError(t *testing.T, method string, url string, data any, code int, msg string) {
	assertions := assert.New(t)
	recorder := performRequest(t, method, url, data)
	checkStatusCode(t, code, recorder.Code)

	if msg != "" {
		response := parseErrorResponse(recorder)
		assertions.Contains(response.Error, msg)
	}
}

func checkFailsWih404(t *testing.T, method string, url string) {
	assertions := assert.New(t)

	recorder := performRequest(t, method, url, nil)
	checkStatusCode(t, http.StatusNotFound, recorder.Code)

	response := parseErrorResponse(recorder)
	assertions.Contains(response.Error, "not found")
}

func deleteSequenceByName(name string) {
	Db.Where("name = ?", name).Delete(&api.Sequence{}) // delete the existing record (if any)
}
