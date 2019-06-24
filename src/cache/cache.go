package cache

import (
	"sync"
	"time"
)

type Item struct {
	Value     interface{}
	Timestamp time.Time
}

func NewItem(v interface{}, t time.Time) *Item {
	return &Item{
		Value:     v,
		Timestamp: t,
	}
}

type Cache struct {
	Lock  sync.RWMutex
	TTL   time.Duration
	Items map[string]*Item
}

func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		TTL:   ttl,
		Items: make(map[string]*Item),
	}
	go func() {
		for {
			time.Sleep(ttl)
			nt := time.Now()
			keys := []string{}
			c.Lock.Lock()
			for k, v := range c.Items {
				t := v.Timestamp
				if nt.Sub(t) > c.TTL {
					keys = append(keys, k)
				}
			}
			for _, k := range keys {
				delete(c.Items, k)
			}
			c.Lock.Unlock()
		}
	}()
	return c
}

func (c *Cache) Get(key string) interface{} {
	c.Lock.RLock()
	var res interface{}
	if it, ok := c.Items[key]; ok {
		res = it.Value
		it.Timestamp = time.Now()
	}
	c.Lock.RUnlock()
	return res
}

func (c *Cache) Put(key string, value interface{}) {
	c.Lock.Lock()
	c.Items[key] = NewItem(value, time.Now())
	c.Lock.Unlock()
}

func (c *Cache) Clear() {
	c.Lock.Lock()
	c.Items = make(map[string]*Item)
}
