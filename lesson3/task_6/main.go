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
		time.Sleep(100 * time.Millisecond) // Имитация задержки записи
	}
	fmt.Printf("Реплика %s закрыта\n", name)
}

func main() {
	wg := sync.WaitGroup{}
	
	input := make(chan int) // Канал для входящих данных
	replicas := []chan int{ // Реплики БД (каналы)
		make(chan int),
		make(chan int),
		make(chan int),
	}

	for i, ch := range replicas {
		wg.Add(1)
		name := fmt.Sprint("Replica: ", i+1) //Перевод формата в строку(норм?)
		go dbReplica(name, ch, &wg)
	}

	go func() {
		for i := 1; i <= 5; i++ {
			input <- i
		}
		close(input)
	}()

	for data := range input {
		var wgTee sync.WaitGroup
		wgTee.Add(len(replicas))

	for _, ch := range replicas {
			go func(c chan int, d int) {
				defer wgTee.Done()
				c <- d
			}(ch, data)
		}

		wgTee.Wait()
	}

	//закрываем каналы всех реплик
	for _, ch := range replicas {
		close(ch) 
	}

	wg.Wait()

	fmt.Println("Все горутины завершены")
}