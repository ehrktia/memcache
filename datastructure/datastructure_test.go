package datastructure

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_get_all(t *testing.T) {
	setup()
	qSize, err := strconv.Atoi(size)
	if err != nil {
		panic("error creating new test cache")
	}
	// test add
	for i := 1; i < qSize; i++ {
		key := fmt.Sprintf("%s-%d", t.Name(), i)
		value := fmt.Sprintf("%d", i)
		_, _ = Add(key, value)
	}
	// test get
	val := GetAll()
	if len(val) != qSize-1 {
		t.Fatalf("expected-%d,got-%d\n", qSize-1, len(val))
	}

}
