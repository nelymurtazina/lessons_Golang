package main

import (
	"sync"
	"time"
)

// func worker() chan int {
// 	ch := make(chan int)
// 	go func() {
// 		time.Sleep(3 * time.Second)
// 		ch <- 42
// 	}()
// 	return ch
// }
// func main() {
// 	timeStart := time.Now()
// 	_, _ = <-worker(), <-worker()
// 	println(int(time.Since(timeStart).Seconds()))
// }

//программа выведет 6 - потому что сначала выполняется 1 горутина (3 сек), потом 2 горутина (это еще + 3 сек)
//нам нужно, чтобы они выполнялись паралельно

func worker(wg *sync.WaitGroup) chan int {
	ch := make(chan int)
	go func() {
		defer wg.Done()
		time.Sleep(3 * time.Second)
		ch <- 42
	}()
	return ch
}

func main() {
	var wg sync.WaitGroup

	timeStart := time.Now()

	for i:=0;i<2;i++{
		wg.Add(1)
		go func ()  {
			_ = <-worker(&wg)
		}()
	}
	wg.Wait()
	println(int(time.Since(timeStart).Seconds()))
}