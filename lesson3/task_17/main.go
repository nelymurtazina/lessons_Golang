package main

import (
	"fmt"
	"log"
	"sync"
)

// Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
// https://victoriametrics.com/blog/go-sync-once/
// https://dev.to/jones_charles_ad50858dbc0/a-developers-guide-to-synconce-your-go-concurrency-lifesaver-3kf2
// https://backendinterview.ru/goLang/concurrency/sync.html

type Plugin interface {
	Execute() string
}

type pluginEntry struct {
	once sync.Once
	plugin Plugin 
	err error 
	initFn func() (Plugin, error) // функция, которая создает плагин
}

type PluginManager struct {
	plugins map[string]*pluginEntry 
	mu sync.RWMutex 
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]*pluginEntry),
	}
}

// RegisterPlugin регистрирует плагин с заданной функцией инициализации
func (pm *PluginManager) RegisterPlugin(name string, initFn func() (Plugin, error)) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.plugins[name] = &pluginEntry{
		initFn: initFn,
	}
}

// GetPlugin возвращает инициализированный плагин по имени
func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	entry, exists := pm.plugins[name]
	pm.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("plugin %s not registered", name)
	}

	// sync.Once гарантирует, что initFn выполнится ровно один раз
	entry.once.Do(func() {
		plugin, err := entry.initFn()
		if err != nil {
			entry.err = err
			return
		}
		entry.plugin = plugin
	})

	// Возвращаем кэшированный результат 
	return entry.plugin, entry.err
}

type DemoPlugin struct{}

func (p *DemoPlugin) Execute() string {
	return "запущен"
}

func initDemo() (Plugin, error) {
	return &DemoPlugin{}, nil
}

func main() {
	pm := NewPluginManager()

	pm.RegisterPlugin("demo", initDemo)

	pm.RegisterPlugin("broken", func() (Plugin, error) {
		return nil, fmt.Errorf("simulated initialization error")
	})

	var wg sync.WaitGroup

	// Тест 1: 5 горутин одновременно запрашивают рабочий плагин
	// Только одна из них выполнит реальную инициализацию
	// Остальные получат уже готовый результат
	fmt.Println("=== Тест рабочего плагина 'demo' ===")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			p, err := pm.GetPlugin("demo")
			if err != nil {
				log.Printf("Горутина %d ошибка: %v", id, err)
				return
			}
			log.Printf("Горутина %d: %s", id, p.Execute())
		}(i)
	}

	wg.Wait()

	// Тест 2: 2 горутины запрашивают плагин с ошибкой инициализации
	// Ошибка кэшируется и возвращается при всех последующих вызовах
	fmt.Println("\n=== Тест плагина с ошибкой 'broken' ===")
	for i := 5; i < 7; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := pm.GetPlugin("broken")
			if err != nil {
				log.Printf("Горутина %d получила ошибку: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	// Тест 3: повторный вызов после всех инициализаций
	// Плагин уже инициализирован, ошибка закэширована
	fmt.Println("\n=== Повторные вызовы ===")
	p, err := pm.GetPlugin("demo")
	if err == nil {
		fmt.Printf("Повторный вызов demo: %s\n", p.Execute())
	}

	_, err = pm.GetPlugin("broken")
	if err != nil {
		fmt.Printf("Повторный вызов broken: %v\n", err)
	}
}