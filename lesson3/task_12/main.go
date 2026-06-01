package main

import (
	"fmt"
	"sync"
	"time"
)

type BoundedQueue struct {
	buffer   []interface{} // срез для хранения элементов
	capacity int           // максимальный размер очереди
	head     int           // индекс головы (откуда забираем)
	tail     int           // индекс хвоста (куда кладём)
	count    int           // текущее количество элементов
	mu       sync.Mutex    // мьютекс для защиты данных
	notFull  *sync.Cond    // сигналит: "очередь НЕ полна"
	notEmpty *sync.Cond    // сигналит: "очередь НЕ пуста"
	closed   bool          // флаг закрытия очереди
}

func NewBoundedQueue(capacity int) *BoundedQueue {
	q := &BoundedQueue{
		buffer:   make([]interface{}, capacity),
		capacity: capacity,
		head:     0,
		tail:     0,
		count:    0,
		closed:   false,
	}
	q.notFull = sync.NewCond(&q.mu)
	q.notEmpty = sync.NewCond(&q.mu)
	return q
}

func (q *BoundedQueue) Put(task interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return
	}

	for q.count == q.capacity {
		q.notFull.Wait()
		if q.closed {
			return
		}
	}

	q.buffer[q.tail] = task
	q.tail = (q.tail + 1) % q.capacity
	q.count++

	q.notEmpty.Signal()
}

func (q *BoundedQueue) Get() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	for q.count == 0 && !q.closed {
		q.notEmpty.Wait()
	}

	if q.count == 0 && q.closed {
		return nil
	}

	item := q.buffer[q.head]
	q.head = (q.head + 1) % q.capacity
	q.count--

	q.notFull.Signal()
	return item
}

func (q *BoundedQueue) Shutdown() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return
	}

	q.closed = true
	q.notFull.Broadcast()
	q.notEmpty.Broadcast()
}

func (q *BoundedQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.count
}

func main() {
	queue := NewBoundedQueue(3)

	go func() {
		for i := 1; i <= 10; i++ {
			queue.Put(fmt.Sprintf("task %d", i))
			time.Sleep(300 * time.Millisecond)
		}
	}()

	go func() {
		for {
			item := queue.Get()
			if item == nil {
				return
			}
			fmt.Printf("processed: %v\n", item)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)
	queue.Shutdown()
	time.Sleep(500 * time.Millisecond)
}