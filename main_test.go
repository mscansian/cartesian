package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeRequest(t *testing.T, method, uri string, expectedCode int) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		t.Errorf("error creating a request: %s", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetPoints)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != expectedCode {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedCode)
	}
	return rr
}

func TestGetPoints(t *testing.T) {
	req := makeRequest(t, http.MethodGet, "/api/points?x=1&y=2&distance=20", http.StatusOK)
	var response interface{}
	if err := json.Unmarshal(req.Body.Bytes(), &response); err != nil {
		t.Errorf("handler returned an  invalid json %s: %s", err, req.Body.String())
	}

}

func TestInvalidRequest(t *testing.T) {
	testTable := []struct {
		url      string
		expected string
	}{
		// Invalid integer
		{"/api/points?x=something&y=5&distance=5", "Must be a valid integer: x"},
		{"/api/points?x=5&y=something&distance=5", "Must be a valid integer: y"},
		{"/api/points?x=5&y=5&distance=something", "Must be a valid integer: distance"},

		//Missing parameters
		{"/api/points", "Missing required parameter: x"},
		{"/api/points?y=5&distance=5", "Missing required parameter: x"},
		{"/api/points?x=5&distance=5", "Missing required parameter: y"},
		{"/api/points?x=5&y=5", "Missing required parameter: distance"},
	}

	for _, tt := range testTable {
		req := makeRequest(t, http.MethodGet, tt.url, http.StatusBadRequest)
		if req.Body.String() != tt.expected {
			t.Errorf("handler returned unexpected body: got %v want %v", req.Body.String(), tt.expected)
		}
	}
}

func TestInvalidMethod(t *testing.T) {
	invalidMethods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}
	for _, method := range invalidMethods {
		t.Logf("testing method %s", method)
		_ = makeRequest(t, method, "/api/points", http.StatusMethodNotAllowed)
	}
}
