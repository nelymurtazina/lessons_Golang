package main

import "fmt"

func Level3() {
    panic("ошибка в Level3")
}

func Level2() {
    defer func() {
        fmt.Println("Завершаем Level2")
    }()
    
    Level3()
}

func Level1() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Паника обработана на уровне 1: %v\n", r)
        }
    }()
    
    Level2()
}

func main() {
    fmt.Println("Начало программы")
    Level1()
    fmt.Println("Программа продолжает работу после обработки паники")
}