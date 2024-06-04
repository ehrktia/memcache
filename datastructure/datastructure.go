package datastructure

import (
	"strings"
	"sync"
)

type Data struct {
	Key   any `json:"key"`
	Value any `json:"value"`
}

// inMemoryCache concurrent cache memory
var inMemoryCache sync.Map
var result = make(chan any)
var load = make(chan bool)
var storeCh = make(chan *Data)

// var errCh = make(chan error)
var done = make(chan bool)

func buildInp(k, v any) {
	go func() {
		d := &Data{
			Key: k, Value: v,
		}
		storeCh <- d
	}()
}

// Add inserts new element in cache.when the element
// is already present it replaces existing value with new
// incoming value and returns new value and loaded flag as true
// triggers update to queue , which adds provided data to end of
// queue
func Add(k, v any) (any, bool) {
	buildInp(k, v)
	go func() {
		data := <-storeCh
		// load data
		_, l := inMemoryCache.LoadOrStore(data.Key, data.Value)
		// return output
		result <- data.Value
		load <- l
		// update queue
		idx := <-inMemoryIdx.getIdx(k)
		// add when missing
		if idx < 0 {
			inMemoryIdx.add(k)
			done <- true
			return
		}
		// when present move to top
		inMemoryIdx.swap(idx)
		done <- true
	}()
	r := <-result
	l := <-load
	return r, l
}

// Delete removes provided key from cache
// remove the key from queue
func Delete(k any) {
	inMemoryCache.Delete(k)
	idx := <-inMemoryIdx.getIdx(k)
	if idx < 0 {
		return
	}
	// remove element from queue
	inMemoryIdx.removeAt(idx)
}

const NotFound = "not found"

func fetch(k any) chan any {
	out := make(chan any)
	go func(key any) {
		v, ok := inMemoryCache.Load(key)
		// when found in cache
		if ok {
			out <- v
			// upd queue
			idx := <-inMemoryIdx.getIdx(k)
			inMemoryIdx.swap(idx)
			return
		}
		out <- NotFound
	}(k)
	return out
}

func Get(k any) any {
	out := <-fetch(k)
	if strings.EqualFold(out.(string), NotFound) {
		return NotFound
	}
	return out
}

// NewQueue creates new inmemory store and queue
// each instance of `once` creates new store and queue
func NewQueue(once *sync.Once, s int) {
	once.Do(func() {
		inMemoryIdx = new(s)
	})
}
