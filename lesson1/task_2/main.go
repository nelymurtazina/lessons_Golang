package main

import (
	"fmt"
	"sort"
	"strings"
)

// WordFrequency принимает строку текста и возвращает map с частотой слов.
func WordFrequency(text string) map[string]int {
	count := make(map[string]int)
	newText := strings.Fields(text)
	for i:= 0; i<len(newText); i++{
		word := newText[i]
		val,ok := count[word]
		if ok {
			count[word]= val + 1
		} else {
			count[word] = 1
		}
	}
	fmt.Println(count)
	PrintWordFrequency(count)
	return count
}

// PrintWordFrequency выводит частотный анализ слов, отсортированный по убыванию частоты.
func PrintWordFrequency(freqMap map[string]int) {
	keys := make([]string, 0, len(freqMap))
	for k := range freqMap{
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	for _, k := range keys {
		fmt.Println(k)
	}
}

func main() {

	text := "golang is great and golang is fast"
	WordFrequency(text)
	
}