package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func worldCount(file string, wg *sync.WaitGroup, results chan <-int) {
	defer wg.Done()
	data, err := os.ReadFile(file)

	if err != nil{
		fmt.Println("Error")
	}
	str := string(data)
	fmt.Println("Text: ", str)

	words := strings.Fields(str)
	count := len(words)
	fmt.Println("Количество слов: ", count)
	
	results <- count
}

func statisticFunc(results <-chan int) {
	sum := 0
	count := 0

	for r := range results {
		sum += r
		count++
	}

	if count<0{
		fmt.Println("Нет данных")
	} else{
		fmt.Println("Всего файлов: ", count)
		fmt.Println("Сумма всех слов: ", sum)
		fmt.Println("Среднее количество слов в файле: ", float64(sum)/float64(count))
	}

}

func main() {
	files := []string{
		"./text/one.txt",
		"./text/two.txt",
		"./text/three.txt",
	}

	results := make(chan int, len(files))
	wg := sync.WaitGroup{}

	for i := 0; i < len(files); i++{
		wg.Add(1)
		go worldCount(files[i], &wg, results)
	}
	
	 go func() {
      wg.Wait()
      close(results)
    }()

		statisticFunc(results)

	
}