package datastructure

import (
	"cmp"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func setup(once *sync.Once) {
	NewQueue(once, 5)
}

func Test_conc_map(t *testing.T) {
	// setup the queue
	once := &sync.Once{}
	setup(once)
	t.Run("add new element to cache", func(t *testing.T) {
		value := fmt.Sprintf("%d", 1)
		result, loaded := Add(t.Name(), value)
		gotValue := reflect.ValueOf(result)
		if value != gotValue.String() {
			t.Fatalf("exp-%s,got-%s", value, gotValue)
		}
		if loaded {
			t.Fatalf("expected-%t,got-%t", false, loaded)
		}
		t.Cleanup(func() {
			Delete(t.Name())
		})
	})
	t.Run("non existing value returns input value and false", func(t *testing.T) {
		nonExistingkey := fmt.Sprintf("%s-%s", t.Name(), t.Name())
		nonExistingvalue := fmt.Sprintf("%s-%s", t.Name(), t.Name())
		r, l := Add(nonExistingkey, nonExistingvalue)
		if l {
			t.Fatalf("expected-%t,got-%t", false, l)
		}
		gotValue := reflect.ValueOf(r)
		if nonExistingvalue != gotValue.String() {
			t.Fatalf("expected-%s,got-%s", nonExistingvalue, gotValue.String())
		}

		t.Cleanup(func() {
			Delete(nonExistingkey)
		})
	})
	t.Run("update existing value return true", func(t *testing.T) {
		// insert data
		value := fmt.Sprintf("%d", 1)
		_, _ = Add(t.Name(), value)
		// update data
		updValue := fmt.Sprintf("%s-%s", value, "2")
		r, l := Add(t.Name(), updValue)
		val := reflect.ValueOf(r)
		if updValue != val.String() {
			t.Fatalf("expected-%s,got-%s", updValue, val.String())
		}
		if !l {
			t.Fatalf("expected-%t,got-%t", true, l)
		}
		t.Cleanup(func() {
			Delete(t.Name())
		})

	})

	t.Run("get existing value from cache", func(t *testing.T) {
		// insert data
		k := t.Name()
		value := fmt.Sprintf("%d", 1)
		_, _ = Add(k, value)
		v, ok := Get(k)
		if !ok {
			t.Fatalf("expected-%t,got-%t", true, ok)
		}
		if v == nil {
			t.Fatalf("expected-%v,got-%v", value, reflect.ValueOf(v).String())
		}
		// check queue position and ensure is top of queue
		val := inMemoryIdx[0]
		if reflect.ValueOf(val).String() != k {
			t.Fatalf("expected-%v,got-%v", k, reflect.ValueOf(val).String())
		}

		t.Cleanup(func() {
			Delete(t.Name())
		})
	})

	t.Run("get non existing value from cache", func(t *testing.T) {
		v, ok := Get(t.Name())
		if ok {
			t.Fatalf("expected-%t,got-%t", false, ok)
		}
		if v != nil {
			t.Fatalf("expected-%v,got-%v", nil, reflect.ValueOf(v).String())
		}
	})
}

func Test_queue(t *testing.T) {
	// setup the queue
	once := &sync.Once{}
	setup(once)
	t.Run("check for space in queue", func(t *testing.T) {
		got := inMemoryIdx.check()
		if cmp.Compare(got, 0) != 0 {
			t.Fatalf("expected-%d\tgot-%d\n", 0, got)
		}
	})

	t.Run("add new element to queue", func(t *testing.T) {
		k := t.Name()
		inMemoryIdx.add(k)
		if inMemoryIdx[0] == nil {
			t.Fatalf("expected-%s,got-%s", t.Name(), inMemoryIdx[0])
		}
	})

	t.Run("evict last element when queue is full", func(t *testing.T) {
		// fill the queue with data
		totLen := len(inMemoryIdx)
		for i := 0; i < totLen-1; i++ {
			k := t.Name()
			inMemoryIdx.add(k)
		}
		// check if queue is full
		for i := 0; i < totLen-1; i++ {
			if inMemoryIdx[i] == nil {
				t.Fatalf("expected-%v,got-%v", t.Name(), inMemoryIdx[i])
			}

		}
		// element to be removed is last element
		got := inMemoryIdx.check()
		if got != totLen-1 {
			t.Fatalf("expected-%d,got-%d", totLen-1, got)
		}
		// evict
		inMemoryIdx.evict()
		if inMemoryIdx[totLen-1] != nil {
			t.Fatalf("expected-%v,got-%v", nil, inMemoryIdx[totLen-1])
		}
	})
}
