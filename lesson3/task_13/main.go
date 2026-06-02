package main

import (
	"fmt"
	"sync"
	"time"
)

type Restaurant struct{
	table []int
	maxTables int
	mu sync.Mutex
	cond  *sync.Cond
}

func NewRestaurant(maxWight int) *Restaurant{
	return &Restaurant{
		table: make([]int, 0, maxWight),
		maxTables: maxWight,
	}
}

func (res *Restaurant) Bronirovanie(wg *sync.WaitGroup){
	defer wg.Done()  // сообщаем, что горутина завершилась
    
	for i := 1; i <= res.maxTables; i++ {
		res.mu.Lock()
		if len(res.table) < res.maxTables {
			res.table = append(res.table, i)
			fmt.Printf("Столик %d заняли. Свободно: %d\n", i, res.maxTables - len(res.table))
		} else {
			fmt.Println("Все столики уже заняты!")
			res.mu.Unlock()
			break
		}
		
		res.mu.Unlock()
		// Имитация времени на обслуживание
		time.Sleep(100 * time.Millisecond)
  }
}


func main(){
	var wg sync.WaitGroup
    
  res := NewRestaurant(5)
    
  wg.Add(1)
  go res.Bronirovanie(&wg)
    
  fmt.Println("Начало бронирования")
  wg.Wait()  // ждем завершения горутины
  fmt.Println("Бронирование завершено")
}