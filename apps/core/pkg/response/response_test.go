package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"lokiforce.com/apps/core/pkg/response"
)

func TestResponse_OK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/ok", func(c *gin.Context) {
		response.OK(c, map[string]string{"foo": "bar"})
	})

	req, _ := http.NewRequest(http.MethodGet, "/ok", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var res response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !res.Success {
		t.Errorf("Expected Success to be true")
	}

	dataMap, ok := res.Data.(map[string]interface{})
	if !ok || dataMap["foo"] != "bar" {
		t.Errorf("Expected Data to contain foo: bar, got %v", res.Data)
	}
}

func TestResponse_Fail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/fail", func(c *gin.Context) {
		response.Fail(c, http.StatusBadRequest, "some error message")
	})

	req, _ := http.NewRequest(http.MethodGet, "/fail", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}

	var res response.Response
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if res.Success {
		t.Errorf("Expected Success to be false")
	}

	if res.Error != "some error message" {
		t.Errorf("Expected Error to be 'some error message', got '%s'", res.Error)
	}
}
