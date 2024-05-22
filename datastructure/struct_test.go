package datastructure

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
)

var size = "5"

func setup(once *sync.Once) {
	os.Setenv("QUEUE_SIZE", size)
	i, err := strconv.Atoi(size)
	if err != nil {
		panic("error creating new test cache")
	}
	NewQueue(once, i)
}

func Test_conc_map(t *testing.T) {
	// setup the queue
	once := &sync.Once{}
	setup(once)
	t.Run("add new element to cache", func(t *testing.T) {
		key := t.Name()
		value := fmt.Sprintf("%d", 1)
		// test add
		r,l:= Add(key, value)
		updIdx := inMemoryIdx.lastAdded
		if updIdx != 0 {
			t.Fatalf("expected-%d\tgot-%d\n", 0, updIdx)
		}
		if r != "1" {
			t.Fatalf("expected:%s\tgot:%s\n", "1", r)
		}
		if l {
			t.Fatalf("expected:%t\tgot:%t\n", false, l)
		}
		if <-done {
			t.Logf("completed")
		}
		if inMemoryIdx.list[0] != t.Name() {
			t.Fatalf("expected-%s\tgot-%s\n", t.Name(), inMemoryIdx.list[0])
		}
		s, err := strconv.Atoi(size)
		if err != nil {
			t.Fatalf("expected:%q\tgot:%v\n", "nil", err)
		}
		if s != len(inMemoryIdx.list) {
			t.Fatalf("expected:%d\tgot-%d\n", s, len(inMemoryIdx.list))
		}
		// test get
		v := Get(key)
		if v.(string) == NotFound {
			t.Fatalf("expected-%q\tgot-%q\n", NotFound, v.(string))
		}
		// test delete
		Delete(t.Name())
		if inMemoryIdx.list[0] != nil {
			t.Fatalf("expected:%q\tgot:%v\n", "nil", inMemoryIdx.list[0])
		}
		if inMemoryIdx.lastAdded != 0 {
			t.Logf("expected:%d\tgot:%d\n", 0, inMemoryIdx.lastAdded)
		}
	})

	t.Run("evict when queue is full", func(t *testing.T) {
		tot := len(inMemoryIdx.list) - 1
		inMemoryIdx.list[tot] = t.Name()
		inMemoryIdx.evict()
		if inMemoryIdx.list[tot] != nil {
			t.Fatalf("expected:%q\tgot:%v\n", "nil", inMemoryIdx.list[tot])
		}
	})

	t.Cleanup(func() {
		err := os.Unsetenv("QUEUE_SIZE")
		if err != nil {
			t.Logf("%v\n", err)
		}
	})
}
