package main

import (
	"fmt"
	"sync"
	"time"
)

type Connection struct {
	ID int
}

type Database struct {
	conn *Connection
	once sync.Once
}

func (db *Database) GetConnection() *Connection {
	db.once.Do(func() {
		fmt.Println("Инициализация подключения к БД")
		time.Sleep(100 * time.Millisecond)
		db.conn = &Connection{ID: 1}
		fmt.Println("Подключение создано!")
	})
	return db.conn
}

func main() {
	db := &Database{}
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			conn := db.GetConnection()
			fmt.Println("Горутина", id, "получила подключение", conn.ID)
		}(i)
	}

	wg.Wait()
}