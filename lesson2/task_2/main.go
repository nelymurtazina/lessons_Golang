package main

import "fmt"

func main() {
	value := 123
	defer fmt.Println(value) //запоминает значение аргумента 123
	changeValue(&value)
}

func changeValue(value *int) {
	*value = 456
}

//Вывод: 123