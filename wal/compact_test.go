package wal

import (
	"fmt"
	"io/fs"
	"os"
	"testing"
)

func Test_compact(t *testing.T) {
	// create a sample wal file locally
	w := new()
	if err := w.createFile(); err != nil {
		t.Fatalf("error-[%v] creating wal file-%q\n", err, w.WalFile())
	}
	// insert some dummy data
	f, err := os.OpenFile(w.WalFile(), os.O_APPEND|os.O_WRONLY, fs.FileMode(os.O_APPEND))
	if err != nil {
		t.Fatalf("error-[%v] opening wal file-%q", err, w.WalFile())
	}
	for i := 1; i <= 3; i++ {
		if _, err := fmt.Fprintf(f, "some test data-%d\n", i); err != nil {
			t.Fatalf("error-[%v] writing data to wal file %q", err, w.WalFile())
		}
	}
	fInfo, err := w.getWalFileInfo()
	if err != nil {
		t.Fatalf("error-[%v] getting wal file info", err)
	}
	// set a default file size
	defaultMaxFileSize = int64(fInfo.Size() - 1)
	t.Logf("file size:%d\n", defaultMaxFileSize)
	if _, err := fmt.Fprintf(f, "some more data to increase size of exisiting file\n"); err != nil {
		t.Fatalf("error-[%v] adding data to trigger compact action", err)
	}
	// test if condition to compact is triggered
	t.Logf("new file size:%d\n", int64(fInfo.Size()))
	if err := f.Close(); err != nil {
		t.Fatalf("error-[%v] closing file %q", err, w.WalFile())
	}
	w.compact()
}
