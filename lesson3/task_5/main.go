package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func readFile(res chan int, file string,wg *sync.WaitGroup){
	defer wg.Done()
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("ОШИБКА")
	}

	str := string(data)
	words := strings.Fields(str)
	count := len(words)

	fmt.Println("Text: ", str)
	fmt.Println("Кол-во слов: ", count)

	res <- count
}

func statistic(ch <- chan int){
	sum := 0
	count := 0

	for r := range ch {
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
	//fan-out!
	files := []string{
		"./text/one.txt",
		"./text/two.txt",
		"./text/three.txt",
	}
	result := make(chan int)

	wg := sync.WaitGroup{}

	for i:=0;i<len(files);i++{
		wg.Add(1)
		go readFile(result, files[i], &wg)
	}

  go func(){
		wg.Wait()
		close(result)
	}()

	
	statistic(result)
}