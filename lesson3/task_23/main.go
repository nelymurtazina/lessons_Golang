package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// https://dev.to/jpoly1219/waitgroups-in-go-3dkj
// https://bytegoblin.io/blog/a-beginners-guide-to-waitgroups-in-go.mdx
// https://www.golinuxcloud.com/golang-waitgroup/
// https://habr.com/ru/articles/850018/

// type logic struct{}

// var Logic logic

// func (l *logic) UpdateDB(ctx context.Context, item *Item) error {
// 	return nil // Заглушка
// }

// func (l *logic) FetchItems(ctx context.Context) ([]*Item, error) {
// 	return []*Item{
// 		{Value: 5},
// 		{Value: 15},
// 		{Value: 7},
// 	}, nil // Заглушка
// }

// type Item struct {
// 	Value int
// }

// func processItem(item *Item) {
// 	time.Sleep(time.Second)
// 	if item.Value > 10 {
// 		fmt.Printf("ERROR: item %d can't be more than 10\n", item.Value)
// 		return
// 	}

// 	err := Logic.UpdateDB(context.Background(), item)
// 	if err != nil {
// 		fmt.Println("ERROR: can't process item")
// 	}
// }

// func DoBusinessLogic() error {
// 	//Нет sync.WaitGroup
// 	items, err := Logic.FetchItems(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	for _, item := range items {
// 		go processItem(item) //запустили горутины и забыли
// 	}

// 	return nil //возвращаем сразу, не дожидаясь горутин
// }

// func main(){
// 	err := DoBusinessLogic()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 	}
// 	fmt.Println("All items processed")
// }


type logic struct{}

var Logic logic

func (l *logic) UpdateDB(ctx context.Context, item *Item) error {
    return nil
}

func (l *logic) FetchItems(ctx context.Context) ([]*Item, error) {
    return []*Item{
        {Value: 5},
        {Value: 15},
        {Value: 7},
    }, nil
}

type Item struct {
    Value int
}

func processItem(item *Item, wg *sync.WaitGroup) {
    defer wg.Done()  // Уменьшаем счётчик при завершении

    time.Sleep(time.Second)
    if item.Value > 10 {
        fmt.Printf("ERROR: item %d can't be more than 10\n", item.Value)
        return
    }

    err := Logic.UpdateDB(context.Background(), item)
    if err != nil {
        fmt.Println("ERROR: can't process item")
    }
}

func DoBusinessLogic() error {
    items, err := Logic.FetchItems(context.Background())
    if err != nil {
        return err
    }

    var wg sync.WaitGroup  // Создаём счётчик горутин

    for _, item := range items {
        wg.Add(1)  // Увеличиваем счётчик
        go processItem(item, &wg)  // Передаём указатель на wg
    }

    wg.Wait()  //Ждём завершения всех горутин
    return nil
}

func main() {
    err := DoBusinessLogic()
    if err != nil {
        fmt.Println("Error:", err)
    }
    fmt.Println("All items processed")
}