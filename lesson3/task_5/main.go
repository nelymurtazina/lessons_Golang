package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

//                     ┌─→ Воркер 1 → результат ─┐
//                     │                         │
// Главная горутина ──→ ├─→ Воркер 2 → результат ─┼→ Сбор результатов
// (раздаёт задачи)    │                         │
//                     └─→ Воркер 3 → результат ─┘

func processerFile(wg *sync.WaitGroup, file string){
	defer wg.Done()

	data, err := os.ReadFile(file)
	if err != nil{
		fmt.Println("Ошибка чтения", err)
	}

	fmt.Println("Text: ",string(data))

	str := string(data)
	words := strings.Fields(str)
	fmt.Println("Количество слов: ", len(words))
}

func main() {
	files := []string{
		"./text/one.txt",
		"./text/two.txt",
		"./text/three.txt",
	}
	wg := sync.WaitGroup{}

	//главная горутина
	for _, file := range files { 
		wg.Add(1)
		go processerFile(&wg, file) 
	}

	wg.Wait()
}