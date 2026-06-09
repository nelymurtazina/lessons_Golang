package main

import (
	"fmt"
	"sync"
	"time"
)

type configManager struct{
	config map[string]string
	once sync.Once
}

func NewConfigManager() *configManager{
	return &configManager{
		config: nil,
	}
}

func (c *configManager) LoadConfig(){
	c.once.Do(func() {
		fmt.Println("Config loaded")
		time.Sleep(100*time.Millisecond)
		c.config = map[string]string{
			"app_name": "MyApp",
			"port":"8080",
			"log_level": "debug",
		}
	fmt.Println("Конфигурация загружена!")
	})
}

func (c *configManager) Get(key string) string{
	c.LoadConfig()
	if val,ok := c.config[key]; ok{
		return val
	}

	return c.config[key] 
}

func (c *configManager) PrintConfig(){
	c.LoadConfig()
	for key, value := range c.config {
		fmt.Printf("%s: %s\n", key, value)
	}
}

func main() {
	cm := &configManager{}
	var wg sync.WaitGroup
	
  for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cm.PrintConfig()
			fmt.Println("Горутина ", id,"app_name = ", cm.Get("app_name"))
		}(i)
	}

	wg.Wait()
}