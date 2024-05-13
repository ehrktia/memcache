package datastructure

type inMemoryQueue []any

var inMemoryIdx inMemoryQueue

// swap replaces first value in queue with
// recently accessed index
func (idx *inMemoryQueue) swap(i int) {
	tmp := (*idx)[0]
	(*idx)[0] = (*idx)[i]
	(*idx)[i] = tmp
}

func (idx *inMemoryQueue) getIndex(k any) int {
	totLen := len(*idx) - 1
	for i := 0; i < totLen; i++ {
		if (*idx)[i] == k {
			return i
		}
	}
	return -1
}

// add inserts the key to last empty position in queue
func (idx *inMemoryQueue) add(k any) {
	index := idx.check()
	(*idx)[index] = k
}

func (idx *inMemoryQueue) check() int {
	// when queue is empty return first position
	if (*idx)[0] == nil {
		return 0
	}
	// get len of queue
	totLen := len(*idx) - 1
	// return the first free space in queue
	for i := totLen; i >= 0; i-- {
		if (*idx)[i] == nil {
			return i
		}
	}
	return totLen

}

// evict follows last accessed/final element
// in queue and remove it from queue and cache store
// this is called only when queue is full and no space a vailable
func (idx *inMemoryQueue) evict() {
	totLen := len(*idx) - 1
	inMemoryIdx[totLen] = nil
}
