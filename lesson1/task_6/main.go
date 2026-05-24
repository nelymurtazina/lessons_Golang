package main

import "fmt"

type account struct {
	value int
}

func main() {
	s1 := make([]account, 0, 2)
	s1 = append(s1, account{})
	s2 := append(s1, account{})
	acc := &s2[0]
	acc.value = 100
	fmt.Println(s1, s2) // s1 = 100, s2=100,0
	s1 = append(s1, account{})
	acc.value += 100
	fmt.Println(s1, s2) //s1 = 200;0 s2 = 200;0
}

// 2.

// func main() {
// 	slice := make([]string, 0, 5)
// 	slice = append(slice, "0", "1", "2", "3")
// 	fmt.Println(slice, len(slice), cap(slice)) //[0 1 2 3] 4, 5
// 	addToSlice1(slice)
// 	fmt.Println(slice, len(slice), cap(slice)) // [0,1,2,one] 4, 5
// 	addToSlice2(slice)
// 	fmt.Println(slice, len(slice), cap(slice)) // [0, 1,2,one] 4, 5! 
// }

// func addToSlice1(slice []string) {
// 	slice = append(slice[1:3], "one")
// }

// func addToSlice2(slice []string) {
// 	slice = append(slice, "two") // [0 1 2 one two] 5 5 
// }

// //3

// func main() {
// 	a1 := make([]int, 0, 10)
// 	a1 = append(a1, []int{1, 2, 3, 4, 5}...) //1,2,3,4,5
// 	a2 := append(a1, 6) // 1,2,3,4,5,7
// 	a3 := append(a1, 7)//1,2,3,4,5,7
// 	fmt.Println(a1, a2, a3) // 
// }

// //4

// func main() {
// 	a := []int{1, 2, 3}
// 	b := a[:2] //указывает на тот же массив, cap хватило 
// 	b = append(b, 4)
// 	fmt.Println(b) // 1,2,4
// 	fmt.Println(a) //1,2,4
// }

// //5

// func main() {
// 	arr := []int{1, 2, 3}
// 	src := arr[:1]
// 	foo(src)
// 	fmt.Println(src) //[1]
// 	fmt.Println(arr) //[1 5 3]
// }

// func foo(src []int) {
// 	src = append(src, 5)
// }

// // 6

// func main() {
// 	arr := [5]int{1, 2, 3, 4, 5}
// 	bar := arr[1:3]
// 	bar = append(bar, 10, 11, 12, 13)
// 	fmt.Println(arr, bar) // [1 2 3 4 5] [2 3 10 11 12 13]
// }

// //7 

// func main() {
// 	a := []string{"a", "b", "c"}
// 	b := a[1:2]
// 	fmt.Println(b, cap(b), len(b)) // b 2 1
// 	b[0] = "q"
// 	fmt.Println(a) // [a q c]
// }

// //8

// func main() {
// 	nums := make([]int, 1, 3)
// 	fmt.Println(nums) // 0
// 	appendSlice(nums, 1)
// 	fmt.Println(nums) // 0
// 	copySlice(nums, []int{2, 3})
// 	fmt.Println(nums) //2
// 	mutateSlice(nums, 1, 4)
// 	fmt.Println(nums) // len 1 -> panic
// }

// func appendSlice(sl []int, val int) {
// 	sl = append(sl, val)
// }

// func copySlice(sl, cp []int) {
// 	copy(sl, cp)
// }

// func mutateSlice(sl []int, idx, val int) {
// 	sl[idx] = val
// }

//9

// func main() {
// 	slice := make([]int, 3, 4)
// 	appendingSlice(slice[:1])
// 	fmt.Println(slice) //
// }

// func appendingSlice(slice []int) {
// 	slice = append(slice, 1)
// }