package main

import (
	"errors"
	"fmt"
)

type MyError struct{
	Code int;
	Msg string
}

func SimpleError() error{
	return errors.New("простая ошибка")
}

func FormattedError(age int) error{
	baseErr:= fmt.Errorf("ошибка: возраст %d недопустим", age)
	return fmt.Errorf("ошибка валидации: %w", baseErr)
}

func (myErr *MyError) Error() string{
	return fmt.Sprintf("код %d: %s", myErr.Code, myErr.Msg)
}

func StructError() error{
	return &MyError{
		Code: 404,
		Msg: "не найдено",
	}
}

func main(){
	err1 := SimpleError()
	fmt.Println("1.", err1)

	err2 := FormattedError(-5)
	fmt.Println("2.", err2)

	err3 := StructError()
	fmt.Println("3.", err3)

	// Проверка доступа к полям структуры
	if myErr, ok := err3.(*MyError); ok {
		fmt.Printf("   Код ошибки: %d, Сообщение: %s\n", myErr.Code, myErr.Msg)
	}
}