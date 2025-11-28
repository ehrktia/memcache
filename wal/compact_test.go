package wal

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"testing"

	"codeberg.org/ehrktia/memcache/datastructure"
)

func Test_write_to_archive(t *testing.T) {
	t.SkipNow()
	if err := os.Setenv(DEFAULT_SIZE, "10"); err != nil {
		t.Fatalf("error-[%v] not able to set default size", err)
	}
	w := NewWal()
	// create file
	if err := CreateFile(w.fileName.fileName); err != nil {
		t.Fatalf("error-[%v] creating file\n", err)
	}
	once := &sync.Once{}
	datastructure.NewQueue(once, 5)
	// write to file
	for i := 1; i <= 5; i++ {
		k := fmt.Sprintf("key-%d", i)
		v := fmt.Sprintf("value-%d", i)
		d := &datastructure.Data{
			Key:   k,
			Value: v,
		}
		b, err := json.Marshal(d)
		if err != nil {
			t.Fatalf("error-[%v] generating data", err)
		}
		if _, err := Write(w.fileName.fileName, b); err != nil {
			t.Fatalf("error-[%v] writing data to file", err)
		}
	}
	if err := Compact(w); err != nil {
		t.Fatalf("error-[%v] writing to archive file", err)
	}
	// add some data to new file
	for i := 1; i <= 2; i++ {
		k := fmt.Sprintf("key-%d", i)
		v := fmt.Sprintf("value-%d", i)
		d := &datastructure.Data{
			Key: k, Value: v,
		}
		db, err := json.Marshal(d)
		if err != nil {
			t.Fatalf("error-[%v] generating data", err)
		}
		if _, err := Write(w.fileName.fileName, db); err != nil {
			t.Fatalf("error-[%v] writing data to file", err)
		}
	}
	t.Cleanup(func() {
		if err := os.Remove(w.fileName.fileName); err != nil {
			t.Fatalf("error-[%v] removing wal file", err)
		}
		if err := os.Remove(w.archiveFileName); err != nil {
			t.Fatalf("error-[%v] removing archive file", err)
		}
		if err := os.Unsetenv(DEFAULT_SIZE); err != nil {
			t.Fatalf("error-[%v] not able to remove defaultsize", err)
		}

	})
}
