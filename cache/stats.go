package cache

import (
	"sync/atomic" // Provides atomic operations for synchronization, for thread-safe counter increments and reads.
	"time"
)

// Stats will hold counters for various operations.
	// uint64 allows for large number of operations and works with atomic package.
type Stats struct {
	Hits		uint64
	Misses		uint64
	Sets		uint64
	Deletes		uint64
	Expirations	uint64
}

// Make a new Stats struct with zero-values
	// Returns a pointer to a new Stats struct
func NewStats() *Stats {
	return &Stats{}
}

// Each of the following methods will use the atomic package to increment the value by 1
	// The use of atomic package here avoids potential raceconditions that could happen with simple increments (i.e. s.Hits++)
		// Simple increment is not atomic and can lead to race conditions
	// Ensures that increments are not lost, even in highly concurrent scenarios.

	// ? An operation is considered "atomic" if it appears to occur instantaneously from the perspective of other threads or goroutines. (The operation is indivisible; either happens completely or not at all, with no observable intermediate states)

func (s *Stats) IncrementHits() 			{ atomic.AddUint64(&s.Hits, 1) }
func (s *Stats) IncrementMisses() 			{ atomic.AddUint64(&s.Misses, 1) }
func (s *Stats) IncrementSets() 			{ atomic.AddUint64(&s.Sets, 1) }
func (s *Stats) IncrementDeletes() 			{ atomic.AddUint64(&s.Deletes, 1) }
func (s *Stats) IncrementExpirations() 		{ atomic.AddUint64(&s.Expirations, 1) }

// Will return a map of Stats for the struct, using atomic package to read from struct
func (s *Stats) GetStats() map[string]uint64 {
	return map[string]uint64{
		"hits":			atomic.LoadUint64(&s.Hits),
		"misses":		atomic.LoadUint64(&s.Misses),
		"sets":			atomic.LoadUint64(&s.Sets),
		"deletes":		atomic.LoadUint64(&s.Deletes),
		"expirations":	atomic.LoadUint64(&s.Expirations),
	}
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