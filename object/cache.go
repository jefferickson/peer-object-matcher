package object

import (
	"sync"
)

// To store cached peer comparisons
type cacheAndMutex struct {
	Cache map[string][]string
	Mu    *sync.RWMutex
}

func (c *cacheAndMutex) get(key string) ([]string, bool) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	temp, ok := c.Cache[key]
	return temp, ok
}

func (c *cacheAndMutex) set(key string, val []string) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	c.Cache[key] = val
}
