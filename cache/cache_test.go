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
