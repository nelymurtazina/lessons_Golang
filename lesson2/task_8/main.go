package main

import (
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

type item struct {
    value     interface{}
    expiresAt time.Time
}

func (i *item) isExpired() bool {
    return time.Now().After(i.expiresAt)
}

type Cache struct {
    mu   sync.RWMutex
    data map[string]*item
}

func NewCache() *Cache {
    return &Cache{
        data: make(map[string]*item),
    }
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data[key] = &item{
        value:     value,
        expiresAt: time.Now().Add(ttl),
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, exists := c.data[key]
    if !exists {
        return nil, false
    }

    if item.isExpired() {
        delete(c.data, key)
        return nil, false
    }

    return item.value, true
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    delete(c.data, key)
}

func (c *Cache) Exists(key string) bool {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, exists := c.data[key]
    if !exists {
        return false
    }

    if item.isExpired() {
        delete(c.data, key)
        return false
    }

    return true
}

func (c *Cache) Clear() {
    c.mu.Lock()
    defer c.mu.Unlock()

    c.data = make(map[string]*item)
}

func (c *Cache) ToJSON() ([]byte, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    snapshot := make(map[string]interface{})

    for key, item := range c.data {
        if !item.isExpired() {
            snapshot[key] = item.value
        }
    }

    return json.Marshal(snapshot)
}

// func (c *Cache) GetAs[T any](key string) (T, error) {
//     var zero T

//     value, ok := c.Get(key)
//     if !ok {
//         return zero, fmt.Errorf("ключ не найден или истек: %s", key)
//     }

//     result, ok := value.(T)
//     if !ok {
//         return zero, fmt.Errorf("тип не соответствует: ожидается %T", zero)
//     }

//     return result, nil
// }

func main() {
    cache := NewCache()

    cache.Set("count", 42, 2*time.Second)
    cache.Set("message", "Hello", 3*time.Second)
    cache.Set("pi", 3.14159, 4*time.Second)

    if value, ok := cache.Get("count"); ok {
        fmt.Println("count:", value)
    }

    fmt.Println("Exists count:", cache.Exists("count"))

    jsonData, _ := cache.ToJSON()
    fmt.Println("JSON:", string(jsonData))

    time.Sleep(3 * time.Second)

    fmt.Println("After 3 seconds:")
    fmt.Println("Exists count:", cache.Exists("count"))
    fmt.Println("Exists message:", cache.Exists("message"))
    fmt.Println("Exists pi:", cache.Exists("pi"))

    // if pi, err := cache.GetAs[float64]("pi"); err == nil {
    //     fmt.Printf("pi = %.5f\n", pi)
    // }

    cache.Clear()
    fmt.Println("After Clear, exists pi:", cache.Exists("pi"))
}