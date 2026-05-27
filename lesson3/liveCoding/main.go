package main

import (
	"fmt"
	"strings"
)

func lengthOfLongestSubstring(s string) int{
		str := string(s)
		res := strings.Split(str, "")
		strRes := make(map[string]bool)
		result := []string{}

		fmt.Println(str, "/", res, len(res))

		for i:= 0 ; i<len(res);i++{
			fmt.Println("Текущая буква:", res[i])
			

			if _, ok := strRes[res[i]]; ok{
				fmt.Println("Найдена повторка! Выводим: ok")
			} else {
				fmt.Println("Новая буква, добавляем в result:", res[i])
				result = append(result, res[i])
			}
			strRes[res[i]] = true 
		}

		
		fmt.Println("Мапа (все увиденные буквы):", strRes)
		fmt.Println("Массив уникальных букв до первой повторки:", result)
		
    return len(result)
}
func main() {
	input := "abcabcbb"
	out:= "aaacdefddd"
	result := lengthOfLongestSubstring(input)
	resultwo := lengthOfLongestSubstring(out)
	fmt.Printf("Длина самой длинной подстроки без повторяющихся символов: %d\n", result) // Вывод: 3
	fmt.Printf("Длина самой длинной подстроки без повторяющихся символов: %d\n", resultwo) // Вывод: 3
}