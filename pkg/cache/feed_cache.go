package cache

import (
	"stargazer/pkg/feed"
	"sync"
	"time"
)

type CacheEntry struct {
	Feed      *feed.RSS
	Timestamp time.Time
}

type FeedCache struct {
	entries map[string]CacheEntry
	mu      sync.RWMutex
}

func NewFeedCache() *FeedCache {
	return &FeedCache{
		entries: make(map[string]CacheEntry),
	}
}

func (c *FeedCache) Get(username string) (*feed.RSS, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[username]
	if !exists {
		return nil, false
	}

	if time.Since(entry.Timestamp) > 15*time.Minute {
		delete(c.entries, username)
		return nil, false
	}

	return entry.Feed, true
}

func (c *FeedCache) Set(username string, feed *feed.RSS) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[username] = CacheEntry{
		Feed:      feed,
		Timestamp: time.Now(),
	}
}
