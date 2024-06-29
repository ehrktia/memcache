// compact check file size
// when size exceeds configured value
// it creates archive and uses a new file
package wal

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"time"
)

var defaultMaxFileSize int64

const DEFAULT_SIZE = "DEFAULT_SIZE"

// getDefaultFileSize gets the configured
// allowed file size for wal, when none defined
// it uses default as size
func getDefaultFileSize() int64 {
	ev := os.Getenv(DEFAULT_SIZE)
	if ev == "" {
		return int64(4096)
	}
	fs, err := strconv.Atoi(ev)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", "error getting wal file size from env var")
		os.Exit(1)
	}
	return int64(fs)
}

// Compact checks the wal file size
// when size exceeds predefined size
// it triggers the compact routine
// creates a new archive file and save it disk
func Compact(w *Wal) error {
	errCh := make(chan error)
	done := make(chan struct{})
	go func() {
		for {
			<-ticker
			// check existing wal file size
			walInfo, err := getWalFileInfo(w.fileName.fileName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error:[%v] getting wal file size", err)
				errCh <- err
			}
			if walInfo.Size() > defaultMaxFileSize {
				fmt.Fprintf(os.Stdout,
					"[info]: start compacting:%s\n", walInfo.Name())
				// create swap file replace wal file
				exisitingFile, err := swapWalFile(w)
				if err != nil {
					errCh <- err
				}
				// create archive file
				archive, err := createArchive(w.archiveFileName)
				if err != nil {
					errCh <- err
				}
				// get file info for existing file
				f, err := os.OpenFile(exisitingFile, os.O_RDONLY, fs.FileMode(os.O_RDONLY))
				if err != nil {
					errCh <- err
				}
				fInfo, err := f.Stat()
				if err != nil {
					errCh <- err
				}
				// archive the wal file
				if err := writeToArchive(archive, f, fInfo); err != nil {
					errCh <- err
				}
				// close the file
				if err := f.Close(); err != nil {
					errCh <- err
				}
				// remove old wal file
				if err := os.Remove(exisitingFile); err != nil {
					errCh <- err
				}
			}
			done <- struct{}{}
		}
	}()
	select {
	case e := <-errCh:
		return e
	case <-done:
		return nil
	}
}

// swapWalFile creates new wal file
func swapWalFile(w *Wal) (string, error) {
	now := time.Now().UnixMicro()
	fname := fmt.Sprintf("wal-%d.txt", now)
	existingFile := w.fileName.fileName
	w.updWalFile(fname)
	if err := CreateFile(w.fileName.fileName); err != nil {
		fmt.Fprintf(os.Stderr,
			"error-[%v] creating new wal file-%q\n",
			err, w.fileName.fileName)
		return existingFile, err
	}
	return existingFile, nil
}

// writeToArchive writes wal file into archive
func writeToArchive(archiveFile, walFile io.ReadWriteCloser,
	walInfo fs.FileInfo) error {
	gzr := gzip.NewWriter(archiveFile)
	defer gzr.Flush()
	defer gzr.Close()
	twr := tar.NewWriter(gzr)
	defer twr.Close()
	tarHeader, err := tar.FileInfoHeader(walInfo, walInfo.Name())
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"error:[%v] generating tar file header for %q",
			err, walInfo.Name())
	}
	tarHeader.Name = walInfo.Name()
	if err := twr.WriteHeader(tarHeader); err != nil {
		fmt.Fprintf(os.Stderr,
			"error:[%v] writing header to tar archive",
			err)
	}
	if _, err := io.Copy(twr, walFile); err != nil {
		fmt.Fprintf(os.Stderr,
			"error:[%v] writing file-%s to archive-%s",
			err, walInfo.Name(), archiveFile)
	}
	return nil
}

func getWalFileInfo(fileName string) (fs.FileInfo, error) {
	f, err := os.Open(fileName)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr,
				"error:[%v] closing wal file-%q\n",
				err, fileName)
		}
	}()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:[%v] opening file\n", err)
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:[%v] getting file stats\n", err)
		return nil, err
	}
	return info, nil
}

// createArchive creates new archive file
func createArchive(fname string) (*os.File, error) {
	archiveFile, err := os.Create(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:[%v] creating archive file with name-%s", err, fname)
		return archiveFile, err
	}
	return archiveFile, err
}
