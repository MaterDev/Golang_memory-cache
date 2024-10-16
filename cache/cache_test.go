package cache

import (
	"testing"
	"time"
)

// Simple unit test to make a new cache item and test that its fields are assigned to the correct values.
func TestCacheItem(t *testing.T) {
	value := "test value"
	expiration := time.Now().Add(time.Hour).UnixNano()

	item := CacheItem{
		Value: value,
		Expiration: expiration,
	}

	if item.Value != value {
		t.Errorf("Expected value %v, got %v", value, item.Value)
	}

	if item.Expiration != expiration {
		t.Errorf("Expected expiration %v, got %v", expiration, item.Expiration)
	}
}

// Will test a setter for the cache for key, value, duration. Duration will be used to calculate expiration, which is some time after now
	// Setter for the 'items' field of cache.
	// If not found, throw error
	// If value doesnt match, throw error.
	// If expiration less than now, throw error (should be in the future)
	// If expiration is greater than expected expiration, throw error (should be in future)
	// Test that stats were made correctly. Will test stats["sets"]
func TestSet(t *testing.T) {
	c := NewCache()
	key := "test_key"
	value := "test_value"
	duration := time.Minute

	c.Set(key, value, duration)

	item, found := c.items[key]

	if !found {
		t.Errorf("Expected item to be stored in cache.")
	}

	if item.Value != value {
		t.Errorf("Expected value %v, got %v", value, item.Value)
	}

	expectedExpiration := time.Now().Add(duration).UnixNano()
	if item.Expiration <= time.Now().UnixNano() || item.Expiration > expectedExpiration {
		t.Errorf("Expected expiration to be in the future, got %v", item.Expiration)
	}

	if c.stats.GetStats()["sets"] != 1 {
		t.Errorf("Expected sets stat to be 1, got %d", c.stats.GetStats()["sets"])
	}
}


// Will test Getter for Cache
func TestGet(t *testing.T) {
	c := NewCache()
	key := "test_key"
	value := "test_value"
	duration := time.Minute

	// Testing to get non-existent key from cache.
	_, found := c.Get(key)

	if found {
		t.Error("Expected key not to be found")
	}
	if c.stats.GetStats()["misses"] != 1 {
		;t.Errorf("Expected misses stat to be 1, got %d", c.stats.GetStats()["misses"])
	}

	// Set key and test get
	c.Set(key, value, duration)
	retrievedValue, found := c.Get(key)
	
	if !found {
		t.Error("Expected key to be found")
	}
	if retrievedValue != value {
		t.Errorf("Expected value %v, got %v", value, retrievedValue)
	}
	if c.stats.GetStats()["hits"] != 1 {
		t.Errorf("Expected hits stat to be 1, got %d", c.stats.GetStats()["hits"])
	}

	// Test expired item (with a negative duration)
	c.Set(key, value, -time.Second)
	_, found = c.Get(key)
	if found {
		t.Error("Expected expired key not to be found")
	}
	// Expect misses to be 2, because it includes the first test assertion for non-existent keys from cache.
	if c.stats.GetStats()["misses"] != 2 {
		t.Errorf("Expected misses stat to be 2, got %d", c.stats.GetStats()["misses"])
	}
}

// Will test removal of a key from c.items. Will increment deletes counter
func TestDelete(t *testing.T) {
	c := NewCache()
	key := "test_key"
	value := "test_value"
	duration := time.Minute
	
	// Set a key
	c.Set(key, value, duration)

	// Delete the key
	c.Delete(key)

	// Try to retrieve the deleted key
	_, found := c.Get(key)
	if found {
		t.Error("Expected key to be deleted")
	}
	if c.stats.GetStats()["deletes"] != 1 {
		t.Errorf("Expected deletes stat to be 1, got %d", c.stats.GetStats()["deletes"])
	}

}