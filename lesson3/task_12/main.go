package main

import (
	"fmt"
	"sync"
	"time"
)

//преимущество sync.Cond — это возможность многократно отправлять сигнал «один ко многим» (Broadcast)
// без пересоздания структуры, сохраняя при этом максимальную производительность.
// В отличие от простого мьютекса, Cond позволяет горутине отпустить блокировку и уснуть,
// а потом проснуться и снова захватить блокировку.

//Вы должны уже захватить мьютекс
// Wait() атомарно: отпускает мьютекс и засыпает
// Когда проснулись — снова захватывает мьютекс
// Проверяем условие снова (поэтому нужен for)

//Продюсеры (кто кладет задачи)
//Консьюмеры (кто забирает задачи)

type BoundedQueue struct {
	mu        sync.Mutex
	isClosed  bool          // переименовано (close - плохое имя)
	maxWeight int           // исправлена опечатка
	queue     []interface{}
	cond      *sync.Cond
}

func NewBoundedQueue(maxWeight int) *BoundedQueue {
	bq := &BoundedQueue{
		queue:     make([]interface{}, 0),
		maxWeight: maxWeight,
		isClosed:  false,
	}
	bq.cond = sync.NewCond(&bq.mu)
	return bq
}

func (bq *BoundedQueue) Put(task interface{}) {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	// Ждем, пока очередь не освободится (НЕ закрыта И есть место)
	for !bq.isClosed && len(bq.queue) >= bq.maxWeight {
		bq.cond.Wait() // Усыпляем продюсера
	}

	// Если очередь закрыта - не добавляем новые задачи
	if bq.isClosed {
		return
	}

	// Добавляем задачу
	bq.queue = append(bq.queue, task)
	
	// Будим ОДНОГО консьюмера (теперь в очереди есть данные)
	bq.cond.Signal()
}

func (bq *BoundedQueue) Get() interface{} {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	// Ждем, пока появятся данные (НЕ закрыта И очередь не пуста)
	for !bq.isClosed && len(bq.queue) == 0 {
		bq.cond.Wait() // Усыпляем консьюмера
	}

	// Если очередь закрыта и данных нет - завершаемся
	if bq.isClosed && len(bq.queue) == 0 {
		return nil
	}

	// Забираем первый элемент
	task := bq.queue[0]
	bq.queue = bq.queue[1:]

	// Будим ОДНОГО продюсера (освободилось место)
	bq.cond.Signal()
	
	return task
}

func (bq *BoundedQueue) Shutdown() {
	bq.mu.Lock()
	defer bq.mu.Unlock()

	bq.isClosed = true
	bq.cond.Broadcast() // Будим ВСЕХ (и продюсеров, и консьюмеров)
}

func main() {
	queue := NewBoundedQueue(3)

	// Продюсеры (5 штук)
	for i := 0; i < 5; i++ {
		go func(id int) {
			for j := 0; j < 10; j++ {
				queue.Put(fmt.Sprintf("task from %d: %d", id, j))
			}
			fmt.Printf("Продюсер %d завершил работу\n", id)
		}(i)
	}

	// Консьюмеры (3 штуки)
	for i := 0; i < 3; i++ {
		go func(id int) {
			for {
				task := queue.Get()
				if task == nil { // Очередь закрыта
					fmt.Printf("Консьюмер %d завершает работу\n", id)
					return
				}
				fmt.Printf("Консьюмер %d обработал: %v\n", id, task)
				time.Sleep(100 * time.Millisecond)
			}
		}(i)
	}

	// Работаем 5 секунд
	time.Sleep(5 * time.Second)
	queue.Shutdown()
	fmt.Println("Очередь закрыта")

	time.Sleep(1 * time.Second) // Даем горутинам завершиться
	fmt.Println("Программа завершена")
}