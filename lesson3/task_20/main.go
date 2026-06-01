package main

import (
    "encoding/json"
    "fmt"
    "sync"
    "time"
)

// Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
// https://ubiklab.net/posts/go-pool-and-mechanics-behind-it/
// https://reliasoftware.com/blog/golang-sync-pool
// https://dev.to/func25/go-syncpool-and-the-mechanics-behind-it-52c1
// https://engineer.yadro.com/article/three-ways-to-optimize-memory-performance-on-go-with-memory-pools/
// https://leapcell.io/blog/boost-go-performance-sync-pool
// https://www.sobyte.net/post/2022-06/go-sync-pool/
// https://goperf.dev/01-common-patterns/object-pooling/

type item struct {
    value      interface{}
    expiration int64
}

type ObjectCache struct {
    items      map[string]item
    mu         sync.RWMutex
    ttl        time.Duration
    bufferPool sync.Pool
}

func NewObjectCache(ttl time.Duration) *ObjectCache {
    c := &ObjectCache{
        items: make(map[string]item),
        ttl:   ttl,
        bufferPool: sync.Pool{
            New: func() interface{} {
                return make([]byte, 0, 4096)
            },
        },
    }
    go c.startCleanup()
    return c
}

func (c *ObjectCache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = item{
        value:      value,
        expiration: time.Now().Add(c.ttl).UnixNano(),
    }
}

func (c *ObjectCache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    item, found := c.items[key]
    if !found {
        return nil, false
    }

    if time.Now().UnixNano() > item.expiration {
        return nil, false
    }

    return item.value, true
}

func (c *ObjectCache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}

func (c *ObjectCache) ToJSON() (string, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    data := make(map[string]interface{})
    now := time.Now().UnixNano()

    for k, v := range c.items {
        if now < v.expiration {
            data[k] = v.value
        }
    }

    jsonBytes, err := json.Marshal(data)
    if err != nil {
        return "", err
    }
    return string(jsonBytes), nil
}

func (c *ObjectCache) startCleanup() {
    ticker := time.NewTicker(c.ttl / 2)
    for range ticker.C {
        c.deleteExpired()
    }
}

func (c *ObjectCache) deleteExpired() {
    c.mu.Lock()
    defer c.mu.Unlock()

    now := time.Now().UnixNano()
    for k, v := range c.items {
        if now > v.expiration {
            delete(c.items, k)
        }
    }
}

func main() {
    cache := NewObjectCache(5 * time.Second)

    cache.Set("user:1", map[string]string{"name": "Alice", "role": "admin"})
    cache.Set("user:2", map[string]string{"name": "Bob", "role": "user"})

    if user, found := cache.Get("user:1"); found {
        fmt.Println("Найден:", user)
    }

    jsonData, _ := cache.ToJSON()
    fmt.Println("Кэш в JSON:", jsonData)

    time.Sleep(6 * time.Second)

    _, found := cache.Get("user:1")
    fmt.Println("После TTL, user:1 найден?", found)
}