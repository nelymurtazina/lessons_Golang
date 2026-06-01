package main

import (
	"fmt"
	"net/http"
	"sync"
	// "time"
)

// func fetchUrl(url string) error {
// 	_, err := http.Get(url)
// 	return err
// }
// func main() {
// 	//Нет sync.WaitGroup
// 	//main не знает, сколько горутин запущено и когда они закончили
// 	urls := []string{
// 		"https://www.lamoda.ru",
// 		"https://www.yandex.ru",
// 		"https://www.mail.ru",
// 		"https://www.google.ru",
// 	}
// 	for _, url := range urls {
// 		go func(url string) {
// 			fmt.Printf("Fetching %s....\n", url)
// 			err := fetchUrl(url)
// 			if err != nil {
// 				fmt.Printf("Error feaching %s: %v\n", url, err)
// 				return
// 			}
// 			fmt.Printf("Fetched %s\n", url)
// 		}(url)
// 	}
// 	fmt.Println("All request launched!")
// 	time.Sleep(400 * time.Millisecond) // ПРОБЛЕМА: ждём фиксированное время вместо синхронизации
// 	//Запросы могут выполняться дольше 400 мс → программа завершится раньше, чем все ответы получены
// 	fmt.Println("Program finished")
// }


func fetchUrl(url string) error {
    _, err := http.Get(url)
    return err
}

func main() {
    urls := []string{
        "https://www.lamoda.ru",
        "https://www.yandex.ru",
        "https://www.mail.ru",
        "https://www.google.ru",
    }

    var wg sync.WaitGroup  // Создаём счётчик горутин

    for _, url := range urls {
        wg.Add(1)  // Увеличиваем счётчик (будет 4)

        go func(url string) {
            defer wg.Done()  //  Уменьшаем счётчик при завершении

            fmt.Printf("Fetching %s...\n", url)
            err := fetchUrl(url)
            if err != nil {
                fmt.Printf("Error fetching %s: %v\n", url, err)
                return
            }
            fmt.Println("Fetched ", url)
        }(url)
    }

    wg.Wait()  // Ждём, пока счётчик не станет 0
    fmt.Println("All requests completed!")
}