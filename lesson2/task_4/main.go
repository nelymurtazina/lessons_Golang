package main

import "fmt"

func handle() error {
	return fmt.Errorf("что-то пошло не так")
}

func main() {
	err := handle()
    if err != nil {
        fmt.Println("Ошибка:", err)
    }

}