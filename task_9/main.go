package main

import (
	"fmt"
	"unsafe"
)

func main() {
	//1
	first := []int{1, 2, 3, 4, 5}
	first = nil
	fmt.Println("first:", first, ":", len(first), ":", cap(first))// first: [] : 0 : 0
	
	//2
	second := []int{1, 2, 3, 4, 5}
	second = second[:0]
	fmt.Println("second:", second, ":", len(second), ":", cap(second)) //second: [] : 0 : 5

	//3
	third := []int{1, 2, 3, 4, 5}
	clear(third)
	fmt.Println("third:", third, ":", len(third), ":", cap(third)) //third: [0 0 0 0 0] : 5 : 5

	//4
	fourth := []int{1, 2, 3, 4, 5}
	clear(fourth[1:3])
	fmt.Println("fourth:", fourth, ":", len(fourth), ":", cap(fourth)) //fourth: [1 0 0 4 5] : 5 : 5

	//5
	slice := make([]int, 3, 6)
	array := [3]int(slice[:3])
	slice[0] = 10

	fmt.Println("slice = ", slice, len(slice), cap(slice)) //slice =  [10 0 0] 3 6
	fmt.Println("array =", array, len(array), cap(array)) //array = [0 0 0] 3 3

	//6 В каких случаях Slice пустой или нулевой
	//1
	var data []string
	fmt.Println("var data []string:")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data)) //empty=true nil=true size=24 data=0x0
	//2
	data = []string(nil)
	fmt.Println("data = []string(nil):")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data)) //empty=true nil=true size=24 data=0x0
	//3
	data = []string{}
	fmt.Println("data = []string{}")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data)) //empty=true nil=false size=24 data=адрес_массива (не 0x0)
	//4
	data = make([]string, 0)
	fmt.Println("data =make([]string,0)")
	fmt.Printf("\tempty=%t nil=%t size=%d data=%p\n", len(data) == 0, data == nil, unsafe.Sizeof(data), unsafe.SliceData(data)) //empty=true nil=false size=24 data=адрес_массива (не 0x0)

	empty := struct{}{}
	fmt.Println("empty struct address ", unsafe.Pointer(&empty)) // empty=true nil=false size=24 data=адрес_массива (не 0x0)
}