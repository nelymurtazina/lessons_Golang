package main

import (
	"fmt"
)

func main() {
	var numbers []*int
	for _, value := range []int{10, 20, 30, 40} {
		//value - одна и та же переменная.
		numbers = append(numbers, &value)
	}
	for _, number := range numbers {
		fmt.Println("d", *number) // 10,20,30,40
	}
}

// //### 2.
// package main

// import (
// 	"fmt"
// 	"strings"
// )

// func chengeSlice(arr []string) {
// 	arr[0] = "Goodbye"
// }

// func appendSomeData(arr []string) string {
// 	arr = append(arr, "!")
// }

// func main() {
// 	someSlice := []string{"Hello", "World"}
// 	chengeSlice(someSlice)
// 	appendSomeData(someSlice)
// 	fmt.Println(strings.Join(someSlice, "")) //GoodbyeWorld
// }
// ```
// ----
// ### 3.
// ```
// package main

// import "fmt"

//всегда ли передается копия заголовков? Путаюсь
//len=2, cap=3
// func test(testSlice []string) {
// 	testSlice = append(testSlice, "Пока")
// }
// func main() {
// 	testSlice := make([]string, 0, 3)
// 	testSlice = append(testSlice, "Привет")
// 	testSlice = append(testSlice, "Привет")
// 	test(testSlice)
// 	fmt.Println(testSlice) //[Привет Привет]
// }
// ```
// ----
// ### 4.
// ```
// package main

// import "fmt"

// func main() {
// 	first := []int{10, 20, 30, 40}
// 	second := make([]*int, len(first))
// 	for i, v := range first {
// 		second[i] = &v
// 	}
// 	fmt.Println(*second[0], *second[1]) //10 20
// }
// ```
// ----
// ### 5.
// ```
// package main

// import (
// 	"fmt"
// )

// func main() {
// 	slice := make([]string, 3, 4)
// 	fmt.Println(slice) //[]

// 	appendSlice(slice)
// 	fmt.Println(slice)//[]

// 	mutareSlice(slice)
// 	fmt.Println(slice)//[vasya]
// }

// func appendSlice(slice []string) {
// 	slice = append(slice, "privet")
// }
// func mutareSlice(slice []string) {
// 	slice[0] = "vasya"
// }