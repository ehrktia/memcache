package wal

import (
	"bytes"
	"encoding/json"
	"os"
	"time"

	"github.com/ehrktia/memcache/datastructure"
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

func (w *Wal) checkHeartBeat() error {
	<-ticker
	// compare file
	isNotSynced, err := w.compareFile()
	if err != nil {
		return err
	}
	if isNotSynced {
		// sync file with in memory
		if err := w.UpdInMemoryCache(); err != nil {
			return err
		}

	}
	return nil

}

func (w *Wal) UpdInMemoryCache() error {
	b := &bytes.Buffer{}
	f, err := os.Open(w.fileName)
	if err != nil {
		return err
	}
	info, err := f.Stat()
	if err != nil {
		return err
	}
	f.Close()
	// when file has data left to read run in loop
	// till all data is read
	for info.Size() > int64(w.filePointer) {
		// read file from last read position
		if err := w.ReadAt(b); err != nil {
			return err
		}
		d := &datastructure.Data{}
		// decode data
		if err := json.NewDecoder(b).Decode(d); err != nil {
			return err
		}
		// upd in memory cache
		_, _ = datastructure.Add(d.Key, d.Value)
		b.Reset()
	}
	return nil
}

// compareFile checks file size with
// file readtill pointer in wal
// returns true when file contains data to be read
func (w *Wal) compareFile() (bool, error) {
	f, err := os.Open(w.fileName)
	if err != nil {
		return false, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return false, err
	}
	if info.Size() > int64(w.filePointer) {
		return true, nil
	}
	return false, nil
}
