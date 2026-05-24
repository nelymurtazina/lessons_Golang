package main

import "fmt"

func main() {  
    fmt.Println("start")  //start
		//end
    for i := 1; i < 4; i++ {  
       defer fmt.Println(i)  
    }  //3 2 1 
    fmt.Println("end") 
}