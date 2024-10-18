package api

import (
	"encoding/json"
	"golang-memory-cache/cache"
	"net/http"
	"time"
)

type Handler struct {
	Cache *cache.Cache
}


func (h *Handler) SetHandler(w http.ResponseWriter, r *http.Request) {
	// If request method is not a post, will error 
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// * url params /set?key=testkey&&value=testvalue&duration=60
	// Will grab data from the URL Query Params
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	durationStr := r.URL.Query().Get("duration")

	if key == "" || value == "" || durationStr == "" {
		http.Error(w, "Missing key, value, or duration", http.StatusBadRequest)
		return
	}

	// Convert the time string into an actual time duration object
	duration, err := time.ParseDuration(durationStr + "s")
	if err != nil {
		http.Error(w, "Invalid duration", http.StatusBadRequest)
	}

	// Will set the cache to hold a new CacheItem with the arguments
	h.Cache.Set(key, value, duration)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	// If request method is not a GET, will  error
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the 'key' query param from the url of the request
	key := r.URL.Query().Get("key")
	// Check if key is provided
	if key == "" {
		http.Error(w, "Missing key parameter", http.StatusBadRequest)
	}

	// Attempt to get the value from the cache
	value, found := h.Cache.Get(key)
	if !found {
		http.Error(w, "Key not found", http.StatusNotFound)
	}

	// Prep response data
		// A map of potential key value pairs, where key is a string and the value is any type
	response := map[string]interface{}{
		"value":value,
	}

	// Set content-type header to application/json (json format)
	w.Header().Set("Content-type", "application/json")

	// encode the response data as JSON and write it to the response writer
		// Encoding is executed during the assignment of the error variable.
		// ! The if statement here is checking if the encoder suceeded basically.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	// If code reaches the end, that means we've successfully written
	// Default 200 OK status
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// If request method is not a DELETE, will  error
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract the 'key' query parameter from URL
	key := r.URL.Query().Get("key")

	// Check that key is provided
	if key == "" {
		http.Error(w, "Missing key, value, or duration", http.StatusBadRequest)
		return
	}

	// Delete the key from the cache
	h.Cache.Delete(key)

	// Write success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Key deleted successfully"))

}
func (h *Handler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
