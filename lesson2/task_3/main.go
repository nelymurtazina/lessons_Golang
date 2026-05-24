package main

import (
	"errors"
	"fmt"
)

func main() {
    // println("Case 1")
    // case1()
    // println()
    // println()

    // println("Case 2")
    // case2()
    // println()
    // println()

    println("Case 3")
    case3()
    println()
    println()

}

func case1() {
    helperWithDefer := func(isError bool) error {
        var retVal error

        defer func() {
            retVal = errors.New("Extra error")
        }()

        if isError {
            retVal = errors.New("Default error")
        }

        return retVal
    }

    helperWithoutDefer := func(isError bool) error {
        var retVal error

        if isError {
            retVal = errors.New("Default error")
        }

        return retVal
    }

    fmt.Println("\twithout:") //without
    fmt.Println(helperWithoutDefer(false))//nil пустой
    fmt.Println(helperWithoutDefer(true))//Default error
    fmt.Println("\twith:")
    fmt.Println(helperWithDefer(false)) // Extra error (у меня вывел компилятор nil, из-за версии?)
    fmt.Println(helperWithDefer(true))//Default error
}

func case2() {
    helperWithDefer := func(isError bool) (retVal error) {
			//Именованное возвращаемое значение retVal существует на протяжении всей функции
        defer func() {
            retVal = errors.New("Extra error")
        }()

        if isError {
            retVal = errors.New("Default error")
        }

        return //выполняются defer и retVal становится "Extra error"
    }

    helperWithoutDefer := func(isError bool) (retVal error) {
        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    fmt.Println("\twithout:")
    fmt.Println(helperWithoutDefer(false)) //nil
    fmt.Println(helperWithoutDefer(true)) //Default error
    fmt.Println("\twith:")
    fmt.Println(helperWithDefer(false)) //Extra error
    fmt.Println(helperWithDefer(true)) //Extra error!
}

func case3() {
    helperWithDefer := func(isError bool) (retVal error) {
        defer func() {
            retVal = errors.New("First Error")
        }()

        defer func() {
            retVal = errors.New("Second Error")
        }()

        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    helperWithoutDefer := func(isError bool) (retVal error) {
        if isError {
            retVal = errors.New("Default error")
        }

        return
    }

    fmt.Println("\twithout:")
    fmt.Println(helperWithoutDefer(false)) //nil
    fmt.Println(helperWithoutDefer(true)) // Default error
    fmt.Println("\twith:")
    fmt.Println(helperWithDefer(false)) // First Error
    fmt.Println(helperWithDefer(true)) // First Error
}