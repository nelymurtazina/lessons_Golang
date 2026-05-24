package main

import (
	"fmt"
	"math"
)

type Shape interface{
	Area() float64
	Perimeter() float64
}

type Circle struct{
	Radius float64
}

type Rectangle struct{
	width float64
	height float64
}

func (c Circle) Area() float64{
	pi := math.Pi
	fmt.Println(c.Radius, "/", 2*pi*c.Radius * c.Radius)
	return  2*pi*c.Radius * c.Radius
	
}
func (r Rectangle) Area() float64{
	fmt.Println(r.width*r.height)
	return r.width*r.height
}

func (c Circle) Perimeter() float64{
	pi := math.Pi
	fmt.Println(2* pi*c.Radius)
	return 2* pi*c.Radius
}

func (r Rectangle) Perimeter() float64{
	fmt.Println(r.width*r.height)
	return 2*(r.width*r.height)
}


func main() {
	fmt.Println("Задание 1")

	circle := Circle{
		Radius: 4.5,
	}
	circle.Area()
}