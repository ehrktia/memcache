package datastructure

import (
	"os"
	"strconv"
	"sync"
	"testing"
)

var size = "5"

func setup() {
	once := &sync.Once{}
	os.Setenv("QUEUE_SIZE", size)
	i, err := strconv.Atoi(size)
	if err != nil {
		panic("error creating new test cache")
	}
	NewQueue(once, i)
}

func Test_conc_map(t *testing.T) {
	// setup the queue
	setup()

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
