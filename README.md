# GoLang Memory Cache

GoLang Memory Cache is an exploratory project implementing a simple, in-memory caching system in Go. This project demonstrates concepts of server-side caching, concurrent-safe operations, and HTTP request handling, without implementing a full server.

![cover](./images/cover.png)

## Project Overview

This project serves as an educational resource and a starting point for understanding how to implement server-side caching in Go. It provides a set of methods for manipulating a cache and includes tests to verify the functionality.

## Features

- In-memory key-value storage with expiration
- Concurrent-safe operations
- Automatic cleanup of expired items
- HTTP handler implementations for cache operations (without an actual server)
- Statistics tracking (hits, misses, sets, deletes, expirations)

## Project Structure

```
golang-memory-cache/
├── cache/
│   ├── cache.go
│   ├── cache_test.go
│   ├── stats.go
│   └── stats_test.go
├── api/
│   ├── handlers.go
│   └── handlers_test.go
```

## Usage

This project is not intended to be used as a standalone application or library. Instead, it serves as a reference implementation and learning tool. You can explore the code, run the tests, and use the concepts demonstrated here in your own projects.

### Exploring the Cache Implementation

The core cache functionality is implemented in `cache/cache.go`. Key methods include:

```go
func NewCache() *Cache
func (c *Cache) Set(key string, value interface{}, duration time.Duration)
func (c *Cache) Get(key string) (interface{}, bool)
func (c *Cache) Delete(key string)
func (c *Cache) GetStats() map[string]uint64
```

### Understanding the Handlers

The `api/handlers.go` file contains handler implementations that demonstrate how the cache might be interacted with via HTTP requests. These handlers are not connected to an actual server but show how such interactions could be structured:

- `SetHandler`: Demonstrates setting a cache item
- `GetHandler`: Demonstrates retrieving a cache item
- `DeleteHandler`: Demonstrates deleting a cache item
- `StatsHandler`: Demonstrates retrieving cache statistics

## Testing

To run the tests and see the cache and handlers in action:

```
go test ./...
```

This will run all tests in the project, demonstrating the functionality of both the cache implementation and the handlers.




