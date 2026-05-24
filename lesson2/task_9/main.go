package main

//1
// import (  
//     "fmt"  
// )  
  
// type MyError struct {  
//     data string  
// }  
  
// func (m *MyError) Error() string {  
//     return m.data  
// }  
// func foo(i int) error {  
//     var err *MyError  
//     if i > 5 {  
//        err = &MyError{data: "i>5"}  
//     }  
//     return err  
// }  
// func main() {  
//     err := foo(4)  
// 		//i>5
//     if err != nil {  
//        fmt.Println("oops")   //"oops"
//     } else {  
//        fmt.Println("ok")  
//     }  
// }



  
// import (  
//     "fmt"  
// )  
  
// type errorString struct {  
//     s string  
// }  
  
// func (e errorString) Error() string {  
//     return e.s  
// }  
  
// func checkErr(err error) {  
//     fmt.Println(err == nil)  
// }  
  
// func main() {  
//     var e1 error //e1 = nil (нет типа, нет значения)
//     checkErr(e1)  //true 
  
//     var e *errorString  
//     checkErr(e)  //false - тк.не равен nil у e есть тип errorString
  
//     e = &errorString{}  
//     checkErr(e)  //false 
  
//     e = nil  //указатель errorString, а значение nil 
//     checkErr(e)  //false 
// }

//3
import "fmt"

type CustomError struct {
	message string
}

func (e *CustomError) Error() string {
	return e.message
}

func returnError(flag bool) error {
	if flag {
		return &CustomError{"Что-то пошло не так"}
	}
	var err *CustomError
	return err
}

func main() {
	err1 := returnError(true) //Что-то пошло не так
	err2 := returnError(false) //nil
	fmt.Println(err2)

	fmt.Println("err1 == nil:", err1 == nil) //false
	fmt.Println("err2 == nil:", err2 == nil) //будет false, тк err2 *CustomError

}
