package main

type Stacker interface {
	Push(v int)
	Pop() int
}

type stack struct {
	data []int
}

func New() *stack {
	return &stack{
		data: make([]int, 0),
	}
}

func (s *stack) Push(v int) {
	s.data = append(s.data, v)
}

func (s *stack) Pop() int {
	if len(s.data) == 0 {
		panic("стек пуст")
	}

	lastIndex := len(s.data) - 1
	lastValue := s.data[lastIndex]
	s.data = s.data[:lastIndex]

	return lastValue
}