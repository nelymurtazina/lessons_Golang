package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func mergeChannels(channels ...<-chan int) <-chan int {
	allCh := make(chan int)

	for _, ch := range channels {
		wg.Add(1)
		go func(abc <-chan int) {
			defer wg.Done()
			for val := range abc {
				allCh <- val
			}
		}(ch)
	}

	go func(){
		wg.Wait()
		close(allCh)
	}()

	return allCh
}

func main() {
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)

	merged := mergeChannels(a, b, c)

    go func() {
        a <- 1
        a <- 2
        a <- 3
        close(a)
    }()

    go func() {
        b <- 4
        b <- 5
        close(b)
    }()

    go func() {
        c <- 6
        close(c)
    }()

    for val := range merged {
       fmt.Println("Val", val)
    }
}