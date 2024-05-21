package datastructure

// type inMemoryQueue []any

type queue struct {
	list      []any
	lastAdded int
}

func new(size int) *queue {
	return &queue{
		list:      make([]any, size),
		lastAdded: 0,
	}
}

func (q *queue) add(k any) {
	q.list[q.lastAdded] = k
	q.lastAdded++
}

func (q *queue) swap(idx int) {
	q.list[0], q.list[idx] = q.list[idx], q.list[0]
}

func (q *queue) getIdx(k any) int {
	totLen := len(q.list) - 1
	for i := 0; i < totLen; i++ {
		if q.list[i] == k {
			return i
		}

	}
	return -1

}

var inMemoryIdx *queue

func (q *queue) evict() {
	lastElementPosition := len(q.list) - 1
	q.list[lastElementPosition] = nil
}

// swap replaces first value in queue with
// recently accessed index
// func (idx *inMemoryQueue) swap(i int) {
// tmp := (*idx)[0]
// (*idx)[0] = (*idx)[i]
// (*idx)[i] = tmp
// }
//
// func (idx *inMemoryQueue) getIndex(k any) int {
// totLen := len(*idx) - 1
// for i := 0; i < totLen; i++ {
// if (*idx)[i] == k {
// return i
// }
// }
// return -1
// }

// add inserts the key to last empty position in queue
// func (idx *inMemoryQueue) add(k any) {
// index := idx.check()
// (*idx)[index] = k
// }

// func (idx *inMemoryQueue) check() int {
// // when queue is empty return first position
// if (*idx)[0] == nil {
// return 0
// }
// // get len of queue
// totLen := len(*idx) - 1
// // return the first free space in queue
// for i := totLen; i >= 0; i-- {
// if (*idx)[i] == nil {
// return i
// }
// }
// return totLen
//
// }

// evict follows last accessed/final element
// in queue and remove it from queue and cache store
// this is called only when queue is full and no space a vailable
// func (idx *) evict() {
// totLen := len(*idx) - 1
// inMemoryIdx[totLen] = nil
// }
