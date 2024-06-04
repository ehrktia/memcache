package wal

import (
	"bytes"
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/ehrktia/memcache/datastructure"
)

func Test_compare_file_size(t *testing.T) {
	// create test file and write some data
	// for testing
	testWal := new()
	fileName := fileName()
	// test file operations
	if err := createFile(); err != nil {
		t.Fatal(err)
	}
	// open file to add some data
	f, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		t.Fatal(err)
	}
	// add data to file
	_, err = f.Write(
		[]byte("this is firstline\nthis is second line\nthis is third line\n"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	// compare pointer to read till
	got, err := testWal.compareFile()
	if err != nil {
		t.Fatal(err)
	}
	// validate
	if !got {
		t.Fatalf("expected-%t\tgot-%t\n", true, got)
	}
	t.Cleanup(func() {
		if err := os.Remove(fileName); err != nil {
			t.Fatal(err)
		}

	})
}

func Test_upd_in_memory_cache(t *testing.T) {
	// create test file and write some data
	// for testing
	testWal := new()
	fileName := fileName()
	// test file operations
	if err := createFile(); err != nil {
		t.Fatal(err)
	}
	// open file to add some data
	f, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		t.Fatal(err)
	}
	d := &datastructure.Data{
		Key: "some-key", Value: "some-value",
	}
	db, err := json.Marshal(&d)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := testWal.Write(db); err != nil {
		t.Fatal(err)
	}
	f.Close()
	b := &bytes.Buffer{}
	// initialize store to test sync
	once := &sync.Once{}
	datastructure.NewQueue(once, 5)
	// test upd cache
	if err := testWal.UpdInMemoryCache(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := os.Remove(fileName); err != nil {
			t.Fatal(err)
		}
		b.Reset()

	})
}
