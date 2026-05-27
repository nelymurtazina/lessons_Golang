package main

import "fmt"

// RemoveUnordered удаляет элемент по индексу без сохранения порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveUnordered[T any](s []T, i int) []T {
	// реализовать
	if i <0 || i >= len(s){
		return s
	}

	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

// RemoveOrdered удаляет элемент по индексу с сохранением порядка.
// Если индекс выходит за границы слайса, возвращает исходный слайс.
func RemoveOrdered[T any](s []T, i int) []T {
	// реализовать

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
	// result := make([]T, 0, len(s))

	//было
	// for _, v := range s {
	// 	if v != value {
	// 		result = append(result, v)
	// 	}
	// }
	// 	return result

	//переделать! переиспользовать RemoveUnordered
	for i := 0; i < len(s); i++ {
        if s[i] == value {
            s = RemoveOrdered(s, i)
            i-- // компенсируем сдвиг из фун-и RemoveOrdered
        }
    }
    return s
	}

// RemoveDuplicates оставляет только уникальные элементы (сохраняет порядок).
func RemoveDuplicates[T comparable](s []T) []T {
	// реализовать
	// if len(s) == 0 {
	// 	return s
	// }
	// seen := make(map[T]bool)

	// result := make([]T, 0, len(s))

	// for _, v := range s {
	// 	if !seen[v] {
	// 		seen[v] = true
	// 		result = append(result, v)
	// 	}
	// }
	// return result

	// через указтели (изменяем исходный слайс)
	//
	if len(s) == 0 {
		return s
	}

	seen := make(map[T]bool)
  writeIndex := 0
	
	for _, v := range s {
    if !seen[v] {
      seen[v] = true
      s[writeIndex] = v  // записываем поверх существующих
      writeIndex++
    }
  }
	return s[:writeIndex]
}

// RemoveIf удаляет элементы, удовлетворяющие условию predicate.
func RemoveIf[T any](s []T, predicate func(T) bool) []T {
	// реализовать
	result := make([]T, 0, len(s))
//1,2,3,4
	for _, v := range s {
		if !predicate(v) {
			result = append(result, v)
			}
		}
		return result

		
}

// RemoveOrderedWithNil удаляет элемент по индексу (для слайса указателей),
// обнуляя удаляемый элемент для предотвращения утечек памяти.
func RemoveOrderedWithNil[T any](s []*T, i int) []*T {
	//реализовать (было)
	// if i < 0 || i >= len(s) {
	// 	return s
	// }
	// s[i] = nil

	// copy(s[i:], s[i+1:])

	// result := s[:len(s)-1]

	// if len(result) < cap(result) {
	// 	result[len(result)-1] = nil
		
	// }
	
	// return result

	//переделать! ошибка nil pointer dereference возникает, 
	// потому что в слайсе после удаления остался nil элемент
	if i<0 || i>= len(s){
		return s
	}
	s[i] = nil
	copy(s[i:], s[i+1:])
	// s[i] = s[len(s)-1]
	result := s[:len(s)-1]
	// result[len(s)-1] = nil
	if cap(result)>len(result){
		result[:cap(result)][len(result)] = nil
	}
	return result
}

// ShrinkCapacity сокращает вместимость слайса, если она превышает
// удвоенную длину после удаления элементов.
func ShrinkCapacity[T any](s []T) []T {
	//реализовать
	if cap(s) <= len(s)*2 {
		return s
	}

	newCap := len(s) + len(s)/10
	if newCap < len(s) {
		newCap = len(s) // на случай, если len=0
	}
	result := make([]T, len(s), newCap)
	copy(result, s)
	
	return result
}


func main() {
	s1 := []int{10, 20, 30, 40, 50}
	fmt.Println("Исходный: ", s1)
	s1 = RemoveUnordered(s1, 2)
	fmt.Println("После удаления индекса: ", s1)

	s2 := []string{"a", "b", "c", "d", "e"}
	fmt.Println("Исходный: ", s2)
	s2 = RemoveOrdered(s2, 2)
	fmt.Println("После удаления индекса 2: ", s2)

	s3 := []int{1, 2, 3, 2, 4, 2, 5}
	fmt.Println("Исходный: ", s3)
	s3 = RemoveAllByValue(s3, 2)
	fmt.Println("После удаления всех 2: ", s3)

	s4 := []string{"яблоко", "банан", "яблоко", "апельсин", "банан", "груша"}
	fmt.Println("Исходный: ", s4)
	s4 = RemoveDuplicates(s4)
	fmt.Println("После удаления дубликатов: ", s4)

	s5 := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("Исходный: ", s5)
	s5 = RemoveIf(s5, func(x int) bool {
		return x%2 == 0
	})
	fmt.Println("После удаления чётных чисел: ", s5)

	a, b, c, d := 10, 20, 30, 40
	s6 := []*int{&a, &b, &c, &d}
	fmt.Println("Исходный:", *s6[0], *s6[1], *s6[2], *s6[3])
	s6 = RemoveOrderedWithNil(s6, 1)
	fmt.Println("После удаления индекса 1: [")
	for i, p := range s6 {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Print(*p)
	}
	fmt.Println("]")

	s7 := make([]int, 3, 100)
	s7[0], s7[1], s7[2] = 1, 2, 3
	fmt.Println("До ShrinkCapacity: ", len(s7), cap(s7), s7)
	s7 = ShrinkCapacity(s7)
	fmt.Println("После ShrinkCapacity: ", len(s7), cap(s7), s7)

	fmt.Println("Работа со строками:")
	words := []string{"hello", "world", "hello", "go", "world", "go", "go"}
	fmt.Println("Исходный: ", words)
	words = RemoveDuplicates(words)
	fmt.Println("Уникальные: ", words)
	
	floats := []float64{1.1, 2.2, 3.3, 2.2, 4.4, 1.1}
	fmt.Println("Исходный: ", floats)
	floats = RemoveAllByValue(floats, 2.2)
	fmt.Println("После удаления 2.2: ", floats)
	
	val1, val2, val3 := 100, 200, 300
	pointers := []*int{&val1, &val2, &val3}
	fmt.Println("До удаления: значения ", *pointers[0], *pointers[1], *pointers[2])
	pointers = RemoveOrderedWithNil(pointers, 1)
	fmt.Println("После удаления: значения ", *pointers[0])
}