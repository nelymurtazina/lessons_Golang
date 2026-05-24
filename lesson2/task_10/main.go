package main

import "fmt"

//1
func CausePanic() {
	panic("Что-то пошло не так")
}

//2
func HandlePanic(){
	defer func() {
		if r := recover(); r != nil{
			fmt.Println("Паника перехвачена: ...")
		}
	}()

	CausePanic()
}

//Обработка паники при делении на ноль

func SafeDivide(a, b int) (result int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Паника перехвачена:", r)
            result = 0
        }
    }()

    if b == 0 {
        panic("деление на ноль")
    }

		fmt.Println(a / b)
    return a / b
}

func main() {
	// `CausePanic()` напрямую
	// CausePanic() //panic потому что вызываем напрямую и панику никто не перехватывает 

	// HandlePanic()
	// fmt.Println("Программа продолжает работу")

	SafeDivide(10, 2) // Ожидаемый результат: 5
  SafeDivide(10, 0) // Ожидаемый результат: 0 (без паники)
}