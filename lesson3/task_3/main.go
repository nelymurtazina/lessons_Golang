package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0 ; i<10;i++{
			number := rand.Intn(10)
			naturals <- number
		}
		close(naturals)
	}()

	wg.Add(1)
	go func(){
		defer wg.Done()
		for val := range naturals{
			fmt.Println("Чтение: ",val)

			squares <- val*val
		}
		close(squares)
	}()

	wg.Add(1)
	go func(){
		defer wg.Done()
		for valTwo := range squares{
			fmt.Println("Квадрат числа", valTwo )
		}
	}()

	wg.Wait()
}