package datastructure

import (
	"fmt"
	"sync"
	"testing"
)

func Test_get_all(t *testing.T) {
	s := 5
	once := &sync.Once{}
	NewQueue(once, s)
	// test add
	for i := 1; i < s; i++ {
		key := fmt.Sprintf("%s-%d", t.Name(), i)
		value := fmt.Sprintf("%d", i)
		r, l := Add(key, value)
		t.Logf("r:%v\n", r)
		t.Logf("loaded:%t\n", l)
	}
	// test get
	val := GetAll()
	for k, v := range val {
		t.Logf("key:%v\tval:%v\n", k, v)
	}
	if len(val) != s-1 {
		t.Fatalf("expected-%d,got-%d\n", s-1, len(val))
	}

}
