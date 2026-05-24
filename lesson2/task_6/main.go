package main

import (
	"errors"
	"fmt"
	"math/rand"
)

var (
       ErrNotFound   = errors.New("ресурс не найден")
       TimeoutError = errors.New("таймаут операции")
)

func SimulateRequest() error{
	num := rand.Intn(11) * 10

	switch {
		case num > 50:
			return fmt.Errorf("запрос не выполнен: %w", TimeoutError)
		case num >30:
			return fmt.Errorf("ошибка: %w", ErrNotFound)
		default:
			return errors.New("неизвестная ошибка")
	}
}

func ProcessError(err error){
	if err == TimeoutError{
		fmt.Println("Требуется повторная попытка")
	} else if err == ErrNotFound{
		fmt.Println("Ресурс не найден")
	} else{
		fmt.Println("Неизвестная ошибка")
	}
}

func main() {

}