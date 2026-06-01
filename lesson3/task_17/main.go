package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
// https://victoriametrics.com/blog/go-sync-once/
// https://dev.to/jones_charles_ad50858dbc0/a-developers-guide-to-synconce-your-go-concurrency-lifesaver-3kf2
// https://backendinterview.ru/goLang/concurrency/sync.html


type Plugin interface {
	Execute() string
}

type pluginEntry struct {
	plugin Plugin
	err    error
}

type PluginManager struct {
	plugins map[string]*pluginEntry
	mu      sync.RWMutex
	initialized map[string]bool
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins:     make(map[string]*pluginEntry),
		initialized: make(map[string]bool),
	}
}

func (pm *PluginManager) RegisterPlugin(name string, initFn func() (Plugin, error)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[name] = &pluginEntry{}
}

func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	entry, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not registered", name)
	}

	if !pm.initialized[name] {
		fmt.Println("Инициализируем плагин ", name)
		time.Sleep(500 * time.Millisecond)

		entry.plugin = &DemoPlugin{}
		entry.err = nil

		pm.initialized[name] = true
		fmt.Println("Плагин инициализирован", name)
	}

	return entry.plugin, entry.err
}

type DemoPlugin struct{}

func (p *DemoPlugin) Execute() string {
	return "DemoPlugin executed!"
}

func main() {
	pm := NewPluginManager()

	pm.RegisterPlugin("demo", func() (Plugin, error) {
		return &DemoPlugin{}, nil
	})
	pm.RegisterPlugin("broken", nil)

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p, err := pm.GetPlugin("demo")
			if err != nil {
				log.Println("Горутина", id, "ошибка:",err)
				return
			}
			log.Println("Горутина", id, p.Execute())
		}(i)
	}

	wg.Wait()
}