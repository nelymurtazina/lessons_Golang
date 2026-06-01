package main

import (
	"fmt"
	"sync"
)

var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 128)
    },
}

func ProcessString(s string) string{
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	if cap(buf) < len(s) {
		buf = make([]byte, len(s))
	} else {
		buf = buf[:len(s)]
	}

	for i, ch := range s {
		if ch >= 'a' && ch <= 'z' {
		buf[i] = byte(ch - 32)
		} else {
			buf[i] = byte(ch)
	}
 	}

	for i, ch := range s{
		if ch >= 'a' && ch <= 'z'{
			buf[i] = byte(ch - 32)
		} else {
			buf[i] = byte(ch)
		}
	}
	return string(buf)
}

func main() {
	examples := []string{
		"hello, world!",
		"gopher",
		"lorem ipsum dolor sit amet",
	}

	for _, s := range examples {
		processed := ProcessString(s)
		fmt.Printf("Original: %q\nProcessed: %q\n\n", s, processed)
	}
}