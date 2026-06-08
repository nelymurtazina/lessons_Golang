package main

import (
	"fmt"
	"sync"
	"unicode"
)

//ПЕРЕДЕЛАТЬ

var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 128)
    },
}

func ProcessString(s string) string{
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	runes := []rune(s)
	neededCap := len(runes)

	if cap(buf) < neededCap {
		buf = make([]byte, 0, neededCap+32)
	} else {
		buf = buf[:0]
	}

	//лучше выделить буфер с запасом! Переделать. Для unicCode не работает.  
	result := make([]byte, 0, len(runes)*4)

	for _, r := range runes {
		// Преобразуем руну в верхний регистр
		upperRune := unicode.ToUpper(r)
		// Добавляем байты руны в результат
		result = append(result, []byte(string(upperRune))...)
	}
	buf = append(buf, result...)
	return string(buf)
}

func main() {
	examples := []string{
		"hello, world!",
		"gopher",
		"lorem ipsum dolor sit amet",
		"Привет, Мир!",
	}

	for _, s := range examples {
		processed := ProcessString(s)
		fmt.Printf("Original: %q\nProcessed: %q\n\n", s, processed)
	}
}