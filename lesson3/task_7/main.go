package main

import (
	"fmt"
	"time"
)

//убрать sync.WaitGroup

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
	go func(){
		for metric := range output{
			fmt.Println("Метрика: ", metric.Name, "Значение: ", metric.Value)
		}
	}()

	input <- ServerMetric{Name: "memory_usage", Value: 8589934592}
	close(input)

	time.Sleep(1* time.Second)
}