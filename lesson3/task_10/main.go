package main

import (
	"fmt"
	"sync"
	"time"
)

//можно безопасно использовать из нескольких горутин (легковесных потоков) одновременно
type SafeCache struct{
	data map[string]string
	mu sync.RWMutex
}

func (c *SafeCache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *SafeCache) Get(key string) (string, bool){
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.data[key]	
	return val, ok
}

func main(){
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

	time.Sleep(time.Second) // wg
	//поработать с неймингом!
}