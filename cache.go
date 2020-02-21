// cache 封装lru支持并发

package basket

import "sync"

// cache 是对lru的并发封装
type cache struct {
	mu       sync.Mutex // 并发锁
	lru      *Lru       // lru
	capacity int        // 容量
}

func (c *cache) Get(key string) (value ByteView, ok bool) {

	c.mu.Lock() // 加锁
	defer c.mu.Unlock()

	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

func (c *cache) Put(key string, value interface{}) {
	c.mu.Lock() // 加锁
	defer c.mu.Unlock()

	if c.lru == nil {
		c.lru = NewLRU(c.capacity)
	}

	c.lru.Put(key, value)
}
