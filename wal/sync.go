package wal

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"
)

// heartBeatInterval used by sync to
// file with in memory cache
var heartBeatInterval = time.Second
var ticker = InitializeTicker()

// check file pointer from wal
// compare file pointer with file size
// when file size is greater than file pointer
// sync file data with in memory cache

func InitializeTicker() <-chan time.Time {
	ticker := time.NewTicker(heartBeatInterval)
	return ticker.C
}

// Write gets data for cache and writes data into
// wal file
func Write(fileName string, data []byte) (int, error) {
	f, err := os.OpenFile(fileName,
		os.O_WRONLY|os.O_APPEND, fs.FileMode(os.O_WRONLY))
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
