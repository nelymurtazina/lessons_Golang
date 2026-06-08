package main

import (
	"fmt"
	"sync"
	"time"
)

// Но если что мы с тобой и так пройдем эти темы. А если хочешь прям догнать,то вот дополнительные ресурсы. Можем отдельно встречу организовать по вопросам::
// https://ubiklab.net/posts/go-sync-cond/
// https://dev.to/func25/go-synccond-the-most-overlooked-sync-mechanism-1fgd
// https://wcademy.ru/go-multithreading-sync-cond/


type Connection struct{
	ID int
}

type ConnectionPool struct{
	connections []*Connection
	maxConnect int
	free []*Connection //свободные соединения
	mu sync.Mutex
	cond *sync.Cond
}

func NewConnectionPool(maxCon int) *ConnectionPool{
	pool := &ConnectionPool{
		free: make([]*Connection, 0, maxCon),
		maxConnect: maxCon,
	}

	// Создаем maxSize соединений
	for i := 1; i <= maxCon; i++{
		conn := &Connection{
			ID:i,
		}
		pool.free = append(pool.free, conn) // все свободны
	}

	pool.cond = sync.NewCond(&pool.mu)
	return pool
}

// Если есть свободное соединение - взять его и вернуть
//Если нет свободных - ждать (cond.Wait)
func (c *ConnectionPool) Get() *Connection{
	c.mu.Lock()
	defer c.mu.Unlock()

	// ПОКА нет свободных И пул не закрыт - жди
	for len(c.free) == 0{
		c.cond.Wait()
	}

	// Берем первое свободное
	conn := c.free[0]
	c.free = c.free[1:]
	return conn
}

func (c *ConnectionPool) Release(conn *Connection){
	c.mu.Lock()
	defer c.mu.Unlock()

	c.free = append(c.free, conn)

	c.cond.Signal()
}


func main() {
    pool := NewConnectionPool(3) // Пул на 3 подключения

    for i := 0; i < 10; i++ {
        go func(id int) {
            conn := pool.Get()
            defer pool.Release(conn)

            fmt.Printf("Горутина %d: подключение %d получено\n", id, conn.ID)
            time.Sleep(2 * time.Second) // Имитация работы
						//контекст
        }(i)
    }

    time.Sleep(10 * time.Second)
}