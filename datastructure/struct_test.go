package datastructure

import (
	"fmt"
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

func TestAdd(t *testing.T) {
	setup()
	got, _ := Add(t.Name(), t.Name())
	if got != t.Name() {
		t.Fatalf("expected testname,got:%s\n", got)
	}
	// if replaced {
	// 	t.Fatal("value not expected to be present")
	// }
}

func BenchmarkAdd(b *testing.B) {
	once := &sync.Once{}
	i := 100
	NewQueue(once, i)
	for j := range i {
		_, _ = Add(b.Name()+fmt.Sprintf("%d", j), b.Name())
	}
}
