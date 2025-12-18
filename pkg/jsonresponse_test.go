package pkg

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResponse_Success(t *testing.T) {
	w := httptest.NewRecorder()
	testData := map[string]string{"key": "value"}

	JSONResponse(w, http.StatusOK, "success", testData)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", contentType)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if response["message"] != "success" {
		t.Errorf("expected message 'success', got %q", response["message"])
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Error("expected data to be a map")
	}
	if data["key"] != "value" {
		t.Errorf("expected data[\"key\"] = 'value', got %v", data["key"])
	}
}

func TestJSONResponse_Error(t *testing.T) {
	w := httptest.NewRecorder()
	errorData := map[string]string{"details": "resource not found"}

	JSONResponse(w, http.StatusNotFound, "not found", errorData)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if response["message"] != "not found" {
		t.Errorf("expected message 'not found', got %q", response["message"])
	}
}

func TestJSONResponse_ComplexData(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"id":      123,
		"name":    "test",
		"active":  true,
		"tags":    []string{"tag1", "tag2"},
		"details": map[string]int{"count": 5},
	}

	JSONResponse(w, http.StatusCreated, "created", data)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	responseData, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected data to be a map")
	}

	if responseData["name"] != "test" {
		t.Errorf("expected name 'test', got %v", responseData["name"])
	}

	if responseData["active"] != true {
		t.Errorf("expected active true, got %v", responseData["active"])
	}
}

func TestJSONResponse_NilData(t *testing.T) {
	w := httptest.NewRecorder()

	JSONResponse(w, http.StatusNoContent, "no content", nil)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status code %d, got %d", http.StatusNoContent, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", contentType)
	}
}
