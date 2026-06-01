package main

import (
	"fmt"
	"sync"
)

type Restaurant struct {
	tables int        // количество свободных столиков
	mu     sync.Mutex // мьютекс для защиты tables
	cond   *sync.Cond // условная переменная
}

func (r *Restaurant) OccupyTable(visitorName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for r.tables == 0 {
		fmt.Println("ждёт свободный столик", visitorName)
		r.cond.Wait()
	}

	r.tables--
	fmt.Println("сел за столик.", visitorName, "Свободно: ", r.tables)
}

func (r *Restaurant) ReleaseTable(visitorName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Освобождаем столик
	r.tables++
	fmt.Println("освободил столик: ", visitorName, "Свободно: ", r.tables)

	// Будим ОДНОГО ждущего (того, кто первым встал в очередь)
	r.cond.Signal()
}