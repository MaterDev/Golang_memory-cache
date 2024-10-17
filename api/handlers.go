package api

import (
	"golang-memory-cache/cache"
	"net/http"
	"time"
)

type Handler struct {
	Cache *cache.Cache
}


func (h *Handler) SetHandler(w http.ResponseWriter, r *http.Request) {
	// If request method is not a post, will return 
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
	// TODO
}
func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
func (h *Handler) StatsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}
