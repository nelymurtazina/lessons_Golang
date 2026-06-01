package main

import (
	"fmt"
	"sync"
)

// Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот ресурсы и пиши по вопросам. Можем отдельно встречу организовать по вопросам:
// https://www.youtube.com/watch?v=luQlkud-jKE&t=5s
// https://habr.com/ru/companies/pt/articles/764850/

func parseFunc(ch <-chan string) <-chan string {
	newCh := make(chan string)
	
	go func(){
		for val := range ch {
		newCh <- val + " parsed - "
		}
		close(newCh)
	}()

	return newCh
}

func split(ch <-chan string, n int) []<-chan string {
	splitChan := make([]chan string, n)

	for i := 0; i < n; i++ {
		splitChan[i] = make(chan string)
	}

	//распределитель 
	go func() {
		defer func() {
		for i := 0; i < n; i++ {
			close(splitChan[i])
		}
	}()

	current := 0

	for data := range ch {
		splitChan[current] <- data
		current = (current + 1) % n //Как работает вообще не поняла, и можно ли использовать переход к другой иттерации без этого?
	}
}()

	result := make([]<-chan string, n)
		for i := 0; i < n; i++ {
		result[i] = splitChan[i]
  }
  return result
	}


func send(ch []<-chan string) <-chan string {
	wg := sync.WaitGroup{}
	newSend := make(chan string)
	
	for _, val := range ch {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for value := range c {
			newSend <- value + "sent - "
			}
		}(val)
	}

	go func(){
		wg.Wait()
		close(newSend)
	}()
	return newSend
}

func main() {
	var wg sync.WaitGroup
	mainCh := make(chan string)
	
	parsed := parseFunc(mainCh)
	splitChannels := split(parsed, 3)
	output := send(splitChannels)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for result := range output {
			fmt.Println("Результат:", result)
		}
	}()

	mainCh <- "One"
	mainCh<- "Two"
	mainCh<- "Three"
	mainCh<- "Four"
	mainCh<- "Five"

	close(mainCh)

	wg.Wait() 
}