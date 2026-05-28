package main

import (
	"fmt"
	"sync"
)

type ServerMetric struct {
	Name  string  // Название метрики (например, "memory_usage")
	Value float64 // Значение в байтах
}

func bytesToMegabytes(ch <-chan ServerMetric) chan ServerMetric {
	chanMegabytes := make(chan ServerMetric  )

	go func(){
		defer close(chanMegabytes)

		for metric := range ch{
			megabytes := metric.Value/(1024*1024)
			chanMegabytes <- ServerMetric{
				Name: metric.Name,
				Value: megabytes,
			}
		}
	}()

	return chanMegabytes
}

func main() {
	input := make(chan ServerMetric)

	output := bytesToMegabytes(input)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(){
		defer wg.Done()
		for metric := range output{
			fmt.Println("Метрика: ", metric.Name, "Значение: ", metric.Value)
		}
	}()

	input <- ServerMetric{Name: "memory_usage", Value: 8589934592}
	close(input)
	wg.Wait()
}