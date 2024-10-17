package api

import (
	"golang-memory-cache/cache"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)


func TestSetHandler(t *testing.T) {
	c := cache.NewCache()
	h := &Handler{Cache: c} // The struct may be very large, so we are going to pass the address instead of assigning, so as to prevent duplicating the data.

	req, err := http.NewRequest("POST", "/set?key=testkey&&value=testvalue&duration=60", nil)
	if err != nil {
		t.Fatal(err)
	}

	// New Response recorder
		// Will act as a mocked version of ResponseWriter
	rr := httptest.NewRecorder()
	//  Converts setHandler into a handler function, which allows regular functions to be executed as HTTP requests
	handler := http.HandlerFunc(h.SetHandler)

	// Call setHandler and pass it the response recorder and request.
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got: %v, expected: %v", status, http.StatusOK)
	}

	// Check if item was sent to the cache.
	value, found := c.Get("testkey")
	if !found {
		t.Error("item was not set in cache")
	}
	if value != "testvalue" {
		t.Errorf("wrong value was set in cache: got %v, expected: %v", value, "testvalue")
	}
}

func TestGetHandler(t *testing.T) {
	// Create new cache and handler
	c := cache.NewCache()
	h := &Handler{Cache: c}

	// Set test item to be in cache
	c.Set("testkey", "testvalue", time.Minute)

	// Create a new GET request with the key as the query param
	req, err := http.NewRequest("GET", "/get?key=testkey", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create the Response Recorder
	rr := httptest.NewRecorder()

	// Create http handler using our GetHandler function
	handler := http.HandlerFunc(h.GetHandler)

	// Serve the request to our handler
	handler.ServeHTTP(rr, req)

	// Check if status code is correct (200 OK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

	// Check if the response body has the correct value
	expected := `{"value":"testvalue}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v, expected %v", rr.Body.String(), expected)
	}

	// Test for non-existent key
	req, _ = http.NewRequest("GET", "/get?key=nonexistent", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check if the status code is 404 , not found for non-existent key
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code for non-existent key: got %v, expected %v", status, http.StatusNotFound)
	}
}