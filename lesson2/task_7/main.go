package main

type Stack[T any] struct{
	elements []T
}

func NewStack[T any]() *Stack[T]{
	return &Stack[T]{
		elements: make([]T, 0),
	}
}

func (s Stack[T]) Push(value T){
	s.elements = append(s.elements, value)
}

func (s Stack[T]) Pop() (T, bool){
	if len(s.elements) == 0{
		var zero T
		return zero, false
	}

	lastIndex := len(s.elements)-1
	value := s.elements[lastIndex]

	s.elements = s.elements[:lastIndex]

	return value,true
}

func (s Stack[T]) Peek() (T, bool){
	if len(s.elements) == 0{
		var zero T
		return zero, false
	}

	lastIndex := len(s.elements)-1
	value := s.elements[lastIndex]
	return value,true
}

func (s Stack[T]) IsEmty() bool{
	return len(s.elements) == 0
}