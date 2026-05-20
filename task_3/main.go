package main

import (
	"fmt"
)


func FilterByValue(m map[int]string, allowedValues []string) map[int]string {
	green := make(map[int]string)
	count := 0

	for i:=1;i<=len(m); i++{
		if m[i] == allowedValues[count]{
			green[i] = m[i]
			count++
		} 
		
	}
	fmt.Println(green)

	return make(map[int]string)
}

// InvertMap меняет ключи и значения местами.
// Если значения исходной map не уникальны, возвращает ошибку.
func InvertMap(m map[string]int) (map[int]string, error) {
	// Проверять уникальность значений
	// При обнаружении дубликата вернуть ошибку с описанием конфликта
	result := make(map[int]string)

	for key, value := range m{
		if _, ok := result[value]; ok{
			return nil, fmt.Errorf("Уже есть такой ключ!")
		} 
		result[value] = key
	}

	fmt.Println(result)
	return result,nil
}

func main() {
	//рандомный пример для проверки
	random := map[int]string {
		1: "Один",
    2: "Два",
    3: "Три",
    4: "Один",
    5: "Пять",
	}

	whiteList := []string{"Один", "Пять"}

	FilterByValue(random, whiteList)

	randomTwo := map[string]int {
    "Анна": 25,
    "Иван": 30,
    "Ольга": 28,
		"Петр": 30,
	}
	//? 
	invertedMap, err := InvertMap(randomTwo)
	if err != nil {
			fmt.Println("Ошибка:", err)
	} else {
			fmt.Println(invertedMap)
	}
}