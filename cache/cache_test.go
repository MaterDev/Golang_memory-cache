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

// Will test method that deletes cached data which has expired
func TestDeleteExpired(t *testing.T) {
	c := NewCache()
	// Expiration will be set as a second before current time.
	c.Set("expired1", "value1", -time.Second)
	c.Set("expired2", "value2", -time.Second)
	// Expiration will be set as a minute after current time.
	c.Set("valid", "value3", time.Minute)

	c.deleteExpired()

	// At the time of the if condition running, we use GET method to attempt retrieval of the data and then run our condition based on whether the data is falsy
	if _, found := c.Get("expired1"); found {
		t.Error("Expected expired1 to be deleted")
	}
	if _, found := c.Get("expired2"); found {
		t.Error("Expected expired2 to be deleted")
	}
	if _, found := c.Get("valid"); !found {
		t.Error("Expected valid to still exist.")
	}

	// Test if expirations has increased for every expired item has been deleted.
	stats := c.stats.GetStats()
	if stats["expirations"] != 2 {
		t.Errorf("Expected expirations stat to be 2, got %d", stats["expirations"])
	}
}

// Te
func TestJanitor(t *testing.T) {
	c := NewCache()
	c.Set("key1", "value1", 2*time.Second)
	c.Set("key2", "value2", 5*time.Second)

	// Wait for the first item to expire, leave the second
	time.Sleep(3 * time.Second)

	// Check if k1 is gone and key2 remains
	if _, found := c.Get("key1"); found {
		t.Error("Expected key1 to be deleted by janitor")
	}
	if _, found := c.Get("key2"); !found {
		t.Error("Expected key2 to remain")
	}

	// Wait additional 3 seconds, for key 2 to 	expire
	time.Sleep(3 * time.Second)

	// Check if key2 is gone
	if _, found := c.Get("key2"); found {
		t.Error("Expected key2 to be deleted by janitor")
	}

	// Stats must be 2
	stats := c.stats.GetStats()
	if stats["expirations"] != 2 {
		t.Errorf("Expected expirations stat to be 2, got %d", stats["expirations"])
	}

	// Stop the janitor
	c.Stop()
}

// Will test that janitor stops when c.Stop() is called
func TestStopJanitor(t *testing.T) {
	c := NewCache()
	c.Set("key", "value", 1*time.Second)

	// Stop the janitor immediately
	c.Stop()

	// Wait 2 seconds for the created item to expire
	time.Sleep(2*time.Second)

	// If item is not there, will error
	if _, found := c.Get("key"); !found {
		t.Error("Expected key to still exist after stopping janitor")
	}

	// Manually try to delete expired item
	c.deleteExpired()

	// Item should still exist because deleteExpired respects the janitor running flag
	if _, found := c.Get("key"); !found {
		t.Error("Expected key to still exist after manual deletion")
	}
}