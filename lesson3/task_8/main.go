package main

import (
	"fmt"
	"time"
)

//убрать sync.WaitGroup

func parsingFunc(ch chan string) chan string {
	outChan := make(chan string)

	go func(){
		for date := range ch{
			time.Sleep(200 * time.Millisecond)
			outChan <- date + " parsed"
		}
		close(outChan)
	}()
	return outChan
}

func RoundRobin(ch <- chan string, n int) []<-chan string {
	outputs := make([]chan string, n)

	for i := 0; i < n; i++ {
		outputs[i] = make(chan string)
	}

	go func() {
		defer func() {
      for i := 0; i < n; i++ {
        close(outputs[i])
      }
	}()

	current := 0 // кому отправляем

	for data := range ch{
		time.Sleep(100 * time.Millisecond)
		outputs[current] <- data
		//оператор взятия остатка от деления
		current = (current + 1) % n
	}
	}()

	result := make([]<-chan string, n)
	for i := 0; i<n;i++{
		result[i] = outputs[i]
	}
	return result
}

func main() {
	chanStr := make(chan string)

	parsStr := parsingFunc(chanStr)

	count := 3
  splitChannels := RoundRobin(parsStr, count)

	for i, ch := range splitChannels {
    go func(id int, inputChan <-chan string) {
      for data := range inputChan {
				time.Sleep(150 * time.Millisecond)
  	    fmt.Println("горутина", id, "получила: ", data)
      }
    }(i, ch)
  }

	go func() {
      chanStr <- "Hello"
      chanStr <- "world"
      chanStr <- "Bye"
      chanStr <- "One"
  	  chanStr <- "two"
      chanStr <- "Golang"
			chanStr <- "Hello"
      chanStr <- "world"
      chanStr <- "Bye"
      chanStr <- "One"
  	  chanStr <- "two"
      chanStr <- "Golang"
      close(chanStr)
  }()
	time.Sleep(10 * time.Second)
}