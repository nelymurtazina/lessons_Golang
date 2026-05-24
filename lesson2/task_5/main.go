package main

import (
	"errors"
	"fmt"
)

type MyError struct{
	Code int;
	Msg string
}

//1
func SimpleError() error {
	return errors.New("простая ошибка")
}

//2
func FormattedError(age int) error{
	baseErr:= fmt.Errorf("ошибка: возраст %d недопустим", age)

	return fmt.Errorf("ошибка валидации: %w", baseErr)
}

//3
func (e MyError) Error() string{
	return e.Msg
}

func StructError() error{
	err := MyError{
		Code: 404,
		Msg: "не найдено",
	}
	return err
}


func main() {

}