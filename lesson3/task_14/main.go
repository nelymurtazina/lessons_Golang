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

type Connection struct {
	ID int  // просто идентификатор подключения
}

type ConnectionPool struct {
	connections []*Connection  // хранилище
	available   int            // сколько свободно сейчас
	capacity    int            // максимальное количество
	mu          sync.Mutex     // защита данных
	cond        *sync.Cond     // условная переменная
}

func NewConnectionPool(capacity int) *ConnectionPool {
	connections := make([]*Connection, capacity)
	for i := 0; i < capacity; i++ {
		connections[i] = &Connection{ID: i + 1}
	}

	p := &ConnectionPool{
		connections: connections,
		available:   capacity,
		capacity:    capacity,
	}
	p.cond = sync.NewCond(&p.mu)
	return p
}

func (p *ConnectionPool) Get() *Connection {
	p.mu.Lock()
	defer p.mu.Unlock()

	for p.available == 0 {
		p.cond.Wait()
	}

	for i, conn := range p.connections {
		if conn != nil {
			p.connections[i] = nil
			p.available--
			return conn
		}
	}
	return nil
}

func (p *ConnectionPool) Release(conn *Connection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := 0; i < p.capacity; i++ {
		if p.connections[i] == nil {
			p.connections[i] = conn
			break
		}
	}
	p.available++

	p.cond.Signal()
}

func main() {
	pool := NewConnectionPool(3)

	for i := 0; i < 10; i++ {
		go func(id int) {
			conn := pool.Get()
			defer pool.Release(conn)

			fmt.Println("Горутина", id, "получила подключение", conn.ID)
			time.Sleep(2 * time.Second)
		}(i)
	}

	time.Sleep(10 * time.Second)
}