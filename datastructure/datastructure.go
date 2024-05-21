package datastructure

// TODO: make it concurrent using channels
import (
	"sync"
)

// inMemoryCache concurrent cache memory
var inMemoryCache sync.Map

// Add inserts new element in cache.when the element
// is already present it replaces existing value with new
// incoming value and returns new value and loaded flag as true
// triggers update to queue , which adds provided data to end of
// queue
func Add(k, v any) (result any, loaded bool) {
	result, loaded = inMemoryCache.LoadOrStore(k, v)
	if loaded {
		result = v
		// update queue with new element
		inMemoryIdx.add(k)
		return
	}
	inMemoryIdx.add(k)
	return
}

// Delete removes provided key from cache
func Delete(k any) {
	inMemoryCache.Delete(k)
}

// Get fetches matching key from cache
// update queue in such manner the recently accessed
// element from cache is moved to front of queue
func Get(k any) (v any, ok bool) {
	v, ok = inMemoryCache.Load(k)
	// when found in cache
	if ok {
		updQueue(k)
	}
	return
}

func updQueue(k any) {
	idx := inMemoryIdx.getIdx(k)
	// when element is missing from queue
	// upd queue and move it to first position
	if idx < 0 {
		lastPosition := len(inMemoryIdx.list) - 1
		inMemoryIdx.list[lastPosition] = k
		inMemoryIdx.swap(lastPosition)

	}

}

// func updQueue(k any) {
// // get element position in queue
// idx := inMemoryIdx.getIndex(k)
// // not found in queue
// if idx < 0 {
// // add new key to queue
// inMemoryIdx.add(k)
// } else {
// // move element to front of queue
// inMemoryIdx.swap(idx)
// }
//
// }

// NewQueue creates new inmemory store and queue
// each instance of `once` creates new store and queue
func NewQueue(once *sync.Once, s int) {
	once.Do(func() {
		inMemoryIdx = new(s)
	})
}
