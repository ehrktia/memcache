package datastructure

import "sync"

type queue struct {
	list      []any
	lastAdded int
	lock      *sync.RWMutex
}

func new(size int) *queue {
	return &queue{
		list:      make([]any, size),
		lastAdded: 0,
		lock:      &sync.RWMutex{},
	}
}

func (q *queue) add(k any) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.list[q.lastAdded] = k
	q.lastAdded++
}

func (q *queue) swap(idx int) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.list[0], q.list[idx] = q.list[idx], q.list[0]
}

func (q *queue) getIdx(k any) chan int {
	out := make(chan int)
	go func() {
		totLen := len(q.list) - 1
		for i := 0; i < totLen; i++ {
			if q.list[i] == k {
				out <- i
			}

		}
		out <- -1
	}()
	return out

}

var inMemoryIdx *queue

func (q *queue) evict() {
	q.lock.Lock()
	defer q.lock.Unlock()
	lastElementPosition := len(q.list) - 1
	q.list[lastElementPosition] = nil
}

func (q *queue) removeAt(idx int) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.list = append(q.list[:idx], q.list[idx+1:]...)
	if q.lastAdded >= 1 {
		q.lastAdded = q.lastAdded - 1
	}

}
