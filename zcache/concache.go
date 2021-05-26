// 并发控制封装
package zcache

import (
	"sync"
	"zcache/lru"
)

type conCache struct {
	cacheBytes int64
	m          sync.Mutex
	lru        *lru.Cache
}

func (c *conCache) add(key string, val ByteView) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes)
	}
	c.lru.Add(key, val)
}

func (c *conCache) get(key string) (val ByteView, ok bool) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
