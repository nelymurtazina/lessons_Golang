package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeCache struct {
	mu   sync.Mutex
	data map[string]string
}

func (cache *SafeCache) Set(key string, value string) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.data[key] = value
}

func (cache *SafeCache) Get(key string) (string, bool){
	cache.mu.Lock()
	defer cache.mu.Unlock()
	val, ok := cache.data[key]
	return val,ok
}

func main() {
	cache := &SafeCache{
		data: make(map[string]string),
	}

	for i := 0; i <100;i++{
		go func(n int){
			key := fmt.Sprintf("key%d", n)
			value := fmt.Sprintf("value%d", n)
			cache.Set(key, value)
		}(i)
	}

		for i := 0; i <100;i++{
		go func(n int){
			key := fmt.Sprintf("key%d", n)
			val,ok := cache.Get(key)
			if ok {
				fmt.Println(val)
			}
		}(i)
	}

	time.Sleep(time.Second) 
}