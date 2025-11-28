package wal

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"codeberg.org/ehrktia/memcache/datastructure"
)

func Test_encode_decode(t *testing.T) {
	testWal := NewWal()
	fileName := testWal.fileName.fileName
	b := &bytes.Buffer{}
	if err := CreateFile(fileName); err != nil {
		t.Fatal(err)
	}
	d := &datastructure.Data{
		Key:   "some-key",
		Value: "some-value",
	}
	if err := json.NewEncoder(b).Encode(d); err != nil {
		t.Fatal(err)
	}
	if _, err := Write(testWal.fileName.fileName, b.Bytes()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Remove(fileName); err != nil {
			t.Fatal(err)
		}
		b.Reset()
	})
}
