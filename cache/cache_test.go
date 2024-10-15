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