package cache
import (
	"testing"
)
// Creation of a new Stats struct
func TestNewStats(t *testing.T) {
	stats := NewStats()

	if stats == nil {
		t.Fatal("Expected NewStats to return a non-nil value")
	}

	if stats.Hits != 0 || stats.Misses != 0 || stats.Sets != 0 || stats.Deletes != 0 || stats.Expirations != 0 {
		t.Error("Expected all stats to be initialized to 0")
	}
}

// Testing function to increment each stat
func TestStatsIncrement(t *testing.T) {
	stats := NewStats()

	stats.IncrementHits()
	if stats.Hits != 1 {
		t.Errorf("Expected Hits to be 1, got %d", stats.Hits)
	}
	stats.IncrementMisses()
	if stats.Misses != 1 {
		t.Errorf("Expected Misses to be 1, got %d", stats.Misses)
	}
	stats.IncrementSets()
	if stats.Sets != 1 {
		t.Errorf("Expected Sets to be 1, got %d", stats.Sets)
	}
	stats.IncrementDeletes()
	if stats.Deletes != 1 {
		t.Errorf("Expected Deletes to be 1, got %d", stats.Deletes)
	}
	stats.IncrementExpirations()
	if stats.Expirations != 1 {
		t.Errorf("Expected Expirations to be 1, got %d", stats.Expirations)
	}
}

func TestGetStats(t *testing.T) {
	stats := NewStats()

	stats.IncrementHits()
	stats.IncrementMisses()
	stats.IncrementSets()
	stats.IncrementDeletes()
	stats.IncrementExpirations()

	result := stats.GetStats()

	expected := map[string]uint64{
		"hits":			1,
		"misses":		1,
		"sets":			1,
		"deletes":		1,
		"expirations":	1,
	}

	for key, value := range expected {
		if result[key] != value {
			t.Errorf("Expected %s to be %d, got %d", key, value, result[key])
		}
	}
}