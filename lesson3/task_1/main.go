package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// func main() {
// 	alreadyStored := make(map[int]struct{})
// 	capacity := 1000
// 	doubles := make([]int, 0, capacity)
// 	// гонка данных. 1000 горутин одновременно читают и пишут в alreadyStored без синхронизации.
// 	//В Go map не безопасен для одновременной записи из нескольких горутин
// 	for i := 0; i < capacity; i++ {
// 		doubles = append(doubles, rand.Intn(10))
// 	}
// 	uniqueIDs := make(chan int, capacity)
// 	wg := sync.WaitGroup{}
// 	for i := 0; i < capacity; i++ {
// 		i := i
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			if _, ok := alreadyStored[doubles[i]]; !ok {
// 				alreadyStored[doubles[i]] = struct{}{}
// 				uniqueIDs <- doubles[i]
// 			}
// 		}()
// 	}

// 	wg.Wait()
// 	//Канал uniqueIDs никогда не закрывается. 
// 	//Цикл for range по каналу будет ждать новые данные вечно → deadlock (вечная блокировка).
// 	for val := range uniqueIDs {
// 		fmt.Println(val)
// 	}
// 	//выведет адрес канала 
// 	fmt.Println(uniqueIDs)
// }

//Исправленный код


func main() {
    // Мьютекс для защиты map (исправление проблемы 1)
    var mu sync.Mutex
    alreadyStored := make(map[int]struct{})
    
    capacity := 1000
    doubles := make([]int, 0, capacity)
    for i := 0; i < capacity; i++ {
        doubles = append(doubles, rand.Intn(10))
    }
    
    uniqueIDs := make(chan int, capacity)
    wg := sync.WaitGroup{}
    
    for i := 0; i < capacity; i++ {
        wg.Add(1)
        go func(idx int) {  // передаём i как параметр (решает проблему захвата)
            defer wg.Done()
            
            mu.Lock()  // блокируем map для безопасного доступа
            _, ok := alreadyStored[doubles[idx]]
            if !ok {
                alreadyStored[doubles[idx]] = struct{}{}
            }
            mu.Unlock()
            
            if !ok {
                uniqueIDs <- doubles[idx]
            }
        }(i)  // передаём текущее значение i
    }
    
    wg.Wait()
    close(uniqueIDs)  // закрываем канал (исправление проблемы 2)
    
    for val := range uniqueIDs {
        fmt.Println(val)
    }
}