package main

import (
	"fmt"
)

// RemoveUnordered удаляет элемент по индексу без сохранения порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveUnordered[T any](s []T, i int) []T {
	// реализовать
	if i < 0 || i >= len(s){
		return  s
	}

	s[i] = s[len(s) - 1]
	return s[:len(s)-1]
}

// RemoveOrdered удаляет элемент по индексу с сохранением порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveOrdered[T any](s []T, i int) []T {
	// реализовать

	if i < 0 || i >= len(s){
		return  s
	}

	for index, v := range s{
		if index == i{
			s = append(s[:i], s[i+1:]...)
			fmt.Println("Удаляем", v)
		}
	}
	return s
}

// RemoveAllByValue удаляет все вхождения указанного значения.
func RemoveAllByValue[T comparable](s []T, value T) []T {
	// реализовать
	result := make([]T, 0, len(s))

	for _, v := range s {
		if v != value {
			result = append(result, v)
		}
	}
	
	return s
}

// RemoveDuplicates оставляет только уникальные элементы (сохраняет порядок).
func RemoveDuplicates[T comparable](s []T) []T {
	// реализовать
	return s
}

// RemoveIf удаляет элементы, удовлетворяющие условию predicate.
func RemoveIf[T any](s []T, predicate func(T) bool) []T {
	// реализовать
	return s
}

// RemoveOrderedWithNil удаляет элемент по индексу (для слайса указателей),
// обнуляя удаляемый элемент для предотвращения утечек памяти.
func RemoveOrderedWithNil[T any](s []*T, i int) []*T {
	//реализовать
	return s
}

// ShrinkCapacity сокращает вместимость слайса, если она превышает
// удвоенную длину после удаления элементов.
func ShrinkCapacity[T any](s []T) []T {
	//реализовать
	return s
}

func main() {
	//реализовать

	ints := []int{10, 20, 30, 40, 50}
	fmt.Println(ints)
	ints = RemoveUnordered(ints, 0)
	fmt.Println(ints)

	ints = RemoveOrdered(ints, 1)
	fmt.Println(ints)
}
