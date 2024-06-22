package wal

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// receive data and persist in file
// keep track of pointer up to where the data is being loaded in to cache

type Wal struct {
	filePointer   int
	size          int
	stamp         int64
	fileName      string
	defaultWalDir string
}

// createFile checks if file exist
// and create new file when not found
func (w *Wal) createFile() error {
	if !w.exists() {
		_, err := os.Create(w.fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// creates file if one does not exist
func new() *Wal {
	now := time.Now().UnixMicro()
	home := os.Getenv("HOME")
	local := ".local"
	return &Wal{
		filePointer:   0,
		size:          4096,
		stamp:         now,
		defaultWalDir: filepath.Join(home, local),
		fileName:      filepath.Join(home, local, fmt.Sprintf("wal-%d.txt", now)),
	}
}

func (w *Wal) WalFile() string {
	return w.fileName
}

func (w *Wal) exists() bool {
	f, err := os.Open(w.fileName)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err := f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "error:[%v] closing-%s file\n", err, w.WalFile())
		os.Exit(1)
	}
	return true
}

var ErrEmptyFile = errors.New("empty file")
var ErrRead = errors.New("no data to read")

// Read uses a temp buffer of size 4096
// read data to buffer and returns number of bytes
// read along with data
func (w *Wal) Read(b *bytes.Buffer) error {
	// open file
	f, err := os.Open(w.fileName)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprint(os.Stderr, "error closing file when trying to Read")
		}
	}()
	if err != nil {
		return err
	}
	// get file size
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	// file empty exit
	if fileInfo.Size() == 0 {
		return ErrEmptyFile
	}
	// reader
	bufReader := bufio.NewReaderSize(f, w.size)
	// write to buffer
	_, err = bufReader.Read(b.Bytes())
	if err != nil && !errors.Is(err, ErrRead) {
		return err
	}
	// when file has no more data
	if errors.Is(err, ErrRead) {
		return ErrRead
	}

	return nil
}

func (w *Wal) Write(data []byte) (int, error) {
	f, err := os.OpenFile(w.fileName, os.O_RDWR, 0666)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprint(os.Stderr, "error closing file when trying to Read")
		}
	}()
	if err != nil {
		return 0, err
	}
	n, err := f.Write(data)
	if err != nil && errors.Is(err, os.ErrPermission) {
		return 0, os.ErrPermission
	}
	if err := f.Sync(); err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Wal) ReadAt(b *bytes.Buffer) error {
	f, err := os.Open(w.fileName)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprint(os.Stderr, "error closing file when trying to Read")
		}
	}()
	if err != nil {
		return err
	}
	buf := make([]byte, w.size)
	n, err := f.ReadAt(buf, int64(w.filePointer))
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	w.filePointer = w.filePointer + n
	// write output to buffer
	_, err = b.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
