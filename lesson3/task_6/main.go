package main

import (
	"fmt"
	"sync"
	"time"
)

func dbReplica(name string, in <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range in {
		fmt.Printf("Запись в %s: %d\n", name, data)
		time.Sleep(500 * time.Millisecond) // Имитация задержки записи
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}


func main() {
	input := make(chan int)
	wg := sync.WaitGroup{}
	var mu sync.Mutex

	replicas := []chan int{
		make(chan int),
		make(chan int),
		make(chan int),
	}

	for i, ch := range replicas {
		wg.Add(1)
		name := fmt.Sprintf("Replica %d", i+1)
		go dbReplica(name, ch, &wg)  // каждая читает из СВОЕГО канала
	}

	go func() {
		for data := range input{
			done := make(chan bool)
			count := 0

			for _, ch := range replicas{
				go func(c chan int,val int){
					c <-val
					mu.Lock()
					count++
					if count == len(replicas){
						done <- true
					}
					mu.Unlock()
				}(ch, data)
			}
			<-done
			close(done)
		}
		for _, ch := range replicas{
			close(ch)
		}
	}()

	for i:=1;i<=5;i++{
		input <-i
	}
	close(input)
	wg.Wait()

}