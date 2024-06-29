package wal

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/ehrktia/memcache/datastructure"
)

func Test_upd_cache(t *testing.T) {
	// create wal file
	w := NewWal()
	if err := CreateFile(w.fileName.fileName); err != nil {
		t.Fatalf("error-[%v] creating file-%q", err, w.fileName.fileName)
	}
	size := 5
	once := &sync.Once{}
	datastructure.NewQueue(once, size)
	d := datastructure.Data{Key: "key", Value: "value"}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("error-[%v] generating test data", err)
	}
	if err := UpdCache(w, b); err != nil {
		t.Fatalf("error-[%v] updating wal and cache with %v\n", err, d)
	}
	v := datastructure.Get("key")
	if strings.EqualFold(v.(string), datastructure.NotFound) {
		t.Fatalf("got-%v\twant-%v\n", v, "value")
	}
	t.Cleanup(func() {
		if err := os.Remove(w.fileName.fileName); err != nil {
			t.Fatalf("error-[%v] removing file", err)
		}

	})

}
