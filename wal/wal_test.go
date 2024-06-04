package wal

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/ehrktia/memcache/datastructure"
)

func Test_encode_decode(t *testing.T) {
	testWal := new()
	fileName := fileName()
	b := &bytes.Buffer{}
	if err := createFile(); err != nil {
		t.Fatal(err)
	}
	d := &datastructure.Data{
		Key:   "some-key",
		Value: "some-value",
	}
	if err := json.NewEncoder(b).Encode(d); err != nil {
		t.Fatal(err)
	}
	if _, err := testWal.Write(b.Bytes()); err != nil {
		t.Fatal(err)
	}
	if err := testWal.Read(b); err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%s\n", b.Bytes())
	t.Cleanup(func() {
		if err := os.Remove(fileName); err != nil {
			t.Fatal(err)
		}
		b.Reset()
	})
}

func Test_encode_read_at(t *testing.T) {
	testWal := new()
	fileName:=fileName()
	b := &bytes.Buffer{}
	if err := createFile(); err != nil {
		t.Fatal(err)
	}
	d := &datastructure.Data{
		Key:   "some-key",
		Value: "some-value",
	}
	db, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := testWal.Write(db); err != nil {
		t.Fatal(err)
	}
	if err := testWal.ReadAt(b); err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%s\n", b.Bytes())
	if len(b.Bytes()) < 1 {
		t.Fatalf("expected:%d\tgot:%d\n", len(b.Bytes()), 100)

	}

	t.Cleanup(func() {
		if err := os.Remove(fileName); err != nil {
			t.Fatal(err)
		}
		b.Reset()
	})

}
