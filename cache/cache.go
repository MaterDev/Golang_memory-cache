package cache

// Represents a single item in the cache
type CacheItem struct {
	Value		interface{}
	Expiration 	int64
}

// The main structure that holds all cached items and statistics
type Cache struct {
	items 	map[string]CacheItem	// Map of key-value pairs of cached items
									// Keys are strings, values are CacheItems
	stats 	*Stats					// Pointer to a Stats object for tracking
}

// Creates and initializes a new Cache instance
func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem), // Initialize an empty map for cache items
		stats: NewStats(), // New Stats object to trck cache operations.
	}
}