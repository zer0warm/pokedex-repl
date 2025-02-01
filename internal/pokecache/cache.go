package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries map[string]cacheEntry

	mutex sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	value     []byte
}

// Create a new cache that will automagically handle entries older than "interval"
func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		Entries: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, value []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Entries[key] = cacheEntry{
		createdAt: time.Now(),
		value:     value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	entry, ok := c.Entries[key]
	return entry.value, ok
}

// Attempt deleting old cache entries every "interval"
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for t := range ticker.C {
		c.reap(t, interval)
	}
}

// Delete cache entries that is older than the latest "interval", from
// "current"
func (c *Cache) reap(current time.Time, interval time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, entry := range c.Entries {
		if entry.createdAt.Add(interval).Before(current) {
			delete(c.Entries, key)
		}
	}
}
