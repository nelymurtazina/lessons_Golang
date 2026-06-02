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
//Может быть попробовать с нуля вместо выполнить это задание. ЧТобы я поняла, как нужно мыслить, когда данно такое условие

type BoundedQueue struct{
	mu sync.Mutex 
	close bool // флаг закрытия
	maxWight int 
	queue []interface{} //как я должна была понять, что тут должен быть интерфейс? из того, что мы передаем в параметр. Чет сложновато((
	cond  *sync.Cond
}

//На что нужно смотреть, чтобы сразу понять, что 100% нужен конструтор? 
// я понимаю для чего он нужен, по типу: создавать и настраивать объект в нач.состоянии.
//КАК НАУЧИТСЯ МЫСЛИТЬ 
func NewBoundedQueue(maxWight int) *BoundedQueue{
	bq := &BoundedQueue{
		queue: make([]interface{}, 0),
		maxWight: maxWight,
		close: false,
	}
	bq.cond = sync.NewCond(&bq.mu) //почему NewCond??
	return bq 
}

func (bq *BoundedQueue) Put(task interface{}){
	bq.mu.Lock()
	defer bq.mu.Unlock()
	// Почему у меня ошибка bq.close == true - почему нельзя явно присвоить ?
	for !bq.close && len(bq.queue) >= bq.maxWight{
		//нужно добавить условие когда очередь заполнена, но я не уверена, что так пишется
		fmt.Println("Ждем чего-то")
		bq.cond.Wait() 
	}

	if bq.close == true{
		return
	}

	fmt.Println("Данные готовы")

	bq.queue = append(bq.queue, task) //добавила элемент
	bq.cond.Signal() // Отправляем сигнал ВСЕМ горутинам

}

func (bq *BoundedQueue) Get() interface{}{
	bq.mu.Lock()
	defer bq.mu.Unlock()
	
	for !bq.close && len(bq.queue) == 0{
		//нужно добавить условие когда очередь заполнена, но я не уверена, что так пишется
		fmt.Println("Очередь пуста! Ждем чего-то")
		bq.cond.Wait() 
	}

	if bq.close == true{
		return nil
	}

	fmt.Println("Данные готовы")

	task := bq.queue[0]
	bq.queue = bq.queue[1:]
	bq.cond.Signal() // Отправляем сигнал ВСЕМ горутинам
	return task
}

func (bq *BoundedQueue) Shutdown(){
	bq.mu.Lock()
	defer bq.mu.Unlock()

	bq.close = true

	//ОКАЗЫВАЕТСЯ! Если есть и спящие продюсеры, и спящие консьюмеры — Signal() разбудит только одного. 
	// Остальные останутся спать навсегда (утечка).
	bq.cond.Broadcast()
}

func main(){
	queue := NewBoundedQueue(3)  // очередь на 3 элемента
    
    // Запускаем 5 продюсеров
    for i := 0; i < 5; i++ {
        go func(id int) {
            for j := 0; j < 10; j++ {
                queue.Put(fmt.Sprintf("task from %d: %d", id, j))
                fmt.Printf("Продюсер %d положил задачу\n", id)
            }
        }(i)
    }
    
    // Запускаем 3 консьюмера
    for i := 0; i < 3; i++ {
        go func(id int) {
            for {
                task := queue.Get()
                if task == nil {  // очередь закрыта
                    fmt.Printf("Консьюмер %d завершает работу\n", id)
                    return
                }
                fmt.Printf("Консьюмер %d обработал: %v\n", id, task)
                time.Sleep(100 * time.Millisecond)  // имитация работы
            }
        }(i)
    }
    
    // Работаем 5 секунд, затем закрываем
    time.Sleep(5 * time.Second)
    queue.Shutdown()
    fmt.Println("Очередь закрыта")
    
    time.Sleep(1 * time.Second)  // даем горутинам завершиться
}