package main

import (
	"fmt"
	"sync"
	"time"
)

type Restaurant struct {
	freeTables int // количество свободных столиков
	totalTables int 
	mu sync.Mutex
	cond *sync.Cond // условная переменная для ожидания
}

// NewRestaurant - создает новый ресторан
func NewRestaurant(totalTables int) *Restaurant {
	r := &Restaurant{
		freeTables:  totalTables, // изначально все столики свободны
		totalTables: totalTables,
	}
	r.cond = sync.NewCond(&r.mu) // связываем cond с мьютексом
	return r
}

// OccupyTable - посетитель занимает столик
// Если свободных столиков нет - горутина засыпает
func (r *Restaurant) OccupyTable(visitorID int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for r.freeTables == 0 {
		fmt.Printf("Посетитель %d: Нет свободных столиков, жду...\n", visitorID)
		r.cond.Wait() // усыпляем горутину (и отпускаем мьютекс)
	}

	r.freeTables--
	fmt.Printf("Посетитель %d занял столик. Свободно столиков: %d/%d\n", 
		visitorID, r.freeTables, r.totalTables)
}

// ReleaseTable - посетитель освобождает столик
func (r *Restaurant) ReleaseTable(visitorID int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.freeTables++
	fmt.Printf("Посетитель %d освободил столик. Свободно столиков: %d/%d\n", 
		visitorID, r.freeTables, r.totalTables)

	// Будим ОДНОГО ожидающего посетителя (Signal)
	r.cond.Signal()
}

func Visitor(id int, restaurant *Restaurant, wg *sync.WaitGroup) {
	defer wg.Done()

	// Занять столик
	restaurant.OccupyTable(id)
	
	time.Sleep(100*time.Millisecond) 
	
	restaurant.ReleaseTable(id)
	
}

func main() {
	restaurant := NewRestaurant(3)
	
	var wg sync.WaitGroup
	totalVisitors := 7 
	
	for i := 1; i <= totalVisitors; i++ {
		wg.Add(1)
		go Visitor(i, restaurant, &wg)
		time.Sleep(300 * time.Millisecond)
	}
	
	wg.Wait()
	
}