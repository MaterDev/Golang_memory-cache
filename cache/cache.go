package cache

import (
	"sync"
	"time"
)

// Represents a single item in the cache
type CacheItem struct {
	Value		interface{}
	Expiration 	int64
}

// The main structure that holds all cached items and statistics
type Cache struct {
	items 	map[string]CacheItem	// Map of key-value pairs of cached items
									// Keys are strings, values are CacheItems
	mu sync.RWMutex
	stats 	*Stats					// Pointer to a Stats object for tracking
	stopJanitor chan bool
	janitorRunning bool
	
}

// Creates and initializes a new Cache instance
func NewCache() *Cache {
	c := &Cache{
		items: make(map[string]CacheItem), // Initialize an empty map for cache items
		stats: NewStats(), // New Stats object to track cache perations.
		stopJanitor: make(chan bool),
		janitorRunning: true,
	}
	go c.janitor()
	return c
}

// Will set
	// expiration, to time.Now() plus incoming duration
	// add a new key to c.items and assign it to a new CacheItem, which has fields to hold incoming value and expiration.
		// ? interface{} type is like "any" type in TS, it is used when the type can be anything.
	// Then will increment stats
	func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
		expiration := time.Now().Add(duration).UnixNano()
	
		c.items[key] = CacheItem{
			Value: value,
			Expiration: expiration,
		}
	
		c.stats.IncrementSets()
	}

// Retireve from cache.items the value for the incoming key and whether it was found. If not found will give (nil, false)
	// when something is not found, will increment misses
	// when something is found, will increment hits
	// If item is expired, will return (nil,false) AND count as a miss
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, found := c.items[key]
	if !found {
		c.stats.IncrementMisses()
		return nil, false
	}

	if c.janitorRunning && time.Now().UnixNano() > item.Expiration {
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

// Will loop through Cache items, if now > expiration, then delete item and increment "Expirations" in stats
func (c *Cache) deleteExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.janitorRunning {
		return // dont delete expired item if janitor is not running.
	}
	now := time.Now().UnixNano()
	for key, item := range c.items {
		if now > item.Expiration {
			delete(c.items, key)
			c.stats.IncrementExpirations()
		}
	}
}

// Will call periodically remove items from the cache that are expired, to free up memory.
func (c *Cache) janitor() {
	ticker := time.NewTicker(time.Second) // will run every second
	defer ticker.Stop() // Will make sure ticker stops at the end of the logic

	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		case <-c.stopJanitor: // Run following sequence of after stopJanitor channel recieves a message.
			c.mu.Lock()
			c.janitorRunning = false
			c.mu.Unlock()
			return
		}
	}
}

func (c *Cache) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.janitorRunning {
		c.stopJanitor <- true
		c.janitorRunning = false
	}
}