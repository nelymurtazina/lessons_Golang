package main

import (
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

type CacheItem struct {
	Value     interface{}
	ExpiresAt time.Time
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]*CacheItem
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]*CacheItem),
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, exists := c.data[key]
	if !exists {
		c.mu.RUnlock()
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.data, key)
		c.mu.Unlock()
		c.mu.RLock()
		return nil, false
	}

	value := item.Value
	c.mu.RUnlock()
	return value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

func (c *Cache) Exists(key string) bool {
	_, ok := c.Get(key)
	return ok
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]*CacheItem)
}

func (c *Cache) ToJSON() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := make(map[string]interface{})
	for key, item := range c.data {
		if !time.Now().After(item.ExpiresAt) {
			snapshot[key] = item.Value
		}
	}
	return json.Marshal(snapshot)
}

func main() {
	cache := NewCache()

	cache.Set("name", "Alice", 5*time.Second)
	cache.Set("age", 25, 2*time.Second)
	cache.Set("score", 98.6, 10*time.Second)

	if val, ok := cache.Get("name"); ok {
		fmt.Println("Имя:", val)
	}

	fmt.Println("Age exists:", cache.Exists("age"))

	jsonBytes, err := cache.ToJSON()
	if err == nil {
		fmt.Println("JSON:", string(jsonBytes))
	}

	time.Sleep(3 * time.Second)

	fmt.Println("Age:", cache.Exists("age"))

	cache.Clear()
	fmt.Println("Name:", cache.Exists("name"))
}