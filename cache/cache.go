package cache

import "time"

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
		stats: NewStats(), // New Stats object to track cache perations.
	}
}

// Retireve from cache.items the value for the incoming key and whether it was found. If not found will give (nil, false)
	// when something is not found, will increment misses
	// when something is found, will increment hits
	// If item is expired, will return (nil,false) AND count as a miss
func (c *Cache) Get(key string) (interface{}, bool) {
	item, found := c.items[key]

	if !found {
		c.stats.IncrementMisses()
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		c.stats.IncrementMisses()
		return nil, false
	}

	c.stats.IncrementHits()
	return item.Value, true
}

// Will delete key from map and increment deletes
func (c *Cache) Delete(key string) {
	// Will delete a specified key from a map
	delete(c.items, key)
	c.stats.IncrementDeletes()
}