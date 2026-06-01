package main

import (
	"fmt"
	"sync"
)

type ConfigManager struct {
	config map[string]string
	once   sync.Once
}

func (cm *ConfigManager) LoadConfig() {
	cm.once.Do(func() {
		fmt.Println("Загрузка конфигурации...")
		cm.config = map[string]string{
			"app_name":  "MyApp",
			"port":      "8080",
			"log_level": "debug",
		}
		fmt.Println("Конфигурация загружена!")
	})
}

func (cm *ConfigManager) Get(key string) string {
	cm.LoadConfig()
	return cm.config[key]
}

func (cm *ConfigManager) PrintConfig() {
	cm.LoadConfig()
	for key, value := range cm.config {
		fmt.Printf("%s: %s\n", key, value)
	}
}

func main() {
	cm := &ConfigManager{}
	var wg sync.WaitGroup

	// Запускаем 10 горутин, которые одновременно пытаются получить доступ к конфигурации
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