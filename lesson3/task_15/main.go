package main

import (
	"fmt"
	"sync"
)

type Connection struct {
	ID int
}

type DataBase struct {
	conn *Connection
	err  error 
	once sync.Once
}

func NewDatabase() *DataBase {
	return &DataBase{}
}

func (db *DataBase) GetConnection() (*Connection, error) {
	db.once.Do(func() {
		// true на false, чтобы протестировать ошибку
		success := true
		
		if success {
			db.conn = &Connection{ID: 1}
			db.err = nil
			fmt.Println("Подключение создано успешно")
		} else {
			db.conn = nil
			db.err = fmt.Errorf("ошибка подключения: сервер БД недоступен")
			fmt.Println(db.err)
		}
	})
	
	return db.conn, db.err
}

func main() {
	db := NewDatabase()
	var wg sync.WaitGroup
	countGo := 10

	for i := 0; i < countGo; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			conn, err := db.GetConnection() 
			if err != nil {
				fmt.Printf("Горутина %d: %v\n", id, err)
				return
			}
			fmt.Printf("Горутина %d: получила соединение %d\n", id, conn.ID)
		}(i)
	}

	wg.Wait()

	conn, err := db.GetConnection()
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		fmt.Printf("Подключение: %d\n", conn.ID)
	}
}