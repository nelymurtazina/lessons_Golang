package main

import (
	"fmt"
	"sync"
	"time"
)

type Connection struct{
	ID int
}

type DataBase struct{
	conn *Connection 
	once sync.Once
}

func (db *DataBase) GetConnection() *Connection{
	//выполнится ТОЛЬКО ОДИН РАЗ
	db.once.Do(func() {
		fmt.Println("Подключение к БД")
		time.Sleep(1 * time.Second)
		db.conn = &Connection{ID: 1}
		fmt.Println("Подключение создано")
	})
	//если нет - создает, сохраняет и возвращает подключение 
	//если есть - возвращают подключение (не создают заново)
	return db.conn
}

func main(){
	db := &DataBase{}

	var wg sync.WaitGroup

	countGo := 10

	for i := 0; i<countGo;i++{
		wg.Add(1)
		go func (id int) {
			defer wg.Done()
			conn := db.GetConnection()
			fmt.Println("Подключение: ", conn.ID)
		}(i)
	}

	wg.Wait()

	fmt.Println("Тест, повторный вызов")
	conn := db.GetConnection()
	fmt.Println("Подключение: ", conn.ID)
}