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
	"path/filepath"
	"strconv"
	"time"
)

var defaultMaxFileSize int64

const DEFAULT_SIZE = "DEFAULT_SIZE"

// getDefaultFileSize gets the configured
// allowed file size for wal, when none defined
// it uses default `3*4096` as size
func getDefaultFileSize() {
	ev := os.Getenv(DEFAULT_SIZE)
	if ev == "" {
		ev = getDefaultSize()
	}
	fs, err := strconv.Atoi(ev)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", "error getting wal file size from env var")
		os.Exit(1)
	}
	defaultMaxFileSize = int64(fs)
}

func getDefaultSize() string {
	dsize := strconv.Itoa(3 * 4096)
	return dsize
}

// compact checks the wal file size
// when size exceeds predefined size `4*4096`
// it triggers the compact routine
// creates a new archive file and save it disk
func (w *Wal) compact() {
	// check existing wal file size
	walInfo, err := w.getWalFileInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error:[%v] getting wal file size", err)
	}
	if walInfo.Size() > defaultMaxFileSize {
		fmt.Fprintf(os.Stdout, "[info]: start compacting:%s\n", walInfo.Name())
		archive, err := createArchive(archiveFileName())
		if err != nil {
			fmt.Fprintf(os.Stderr,
				"error:[%v] creating archive with name-%s\n", err, archiveFileName())
		}
		if err := w.writeToArchive(archive, walInfo, archiveFileName()); err != nil {
			fmt.Fprintf(os.Stderr, "error:[%v] writing file to archive", err)
		}
	}
}

func (w *Wal) writeToArchive(archiveFile io.ReadWriter,
	walInfo fs.FileInfo, archiveFileName string) error {
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
			"error:[%v] writing header to tar archive-%s",
			err, archiveFileName)
	}
	f, err := os.OpenFile(w.fileName, os.O_RDONLY, fs.FileMode(os.O_RDONLY))
	if err != nil {
		fmt.Fprintf(os.Stderr,
			"error:[%v] opening wal file-%s to copy into archive",
			err, walInfo.Name())
	}
	if _, err := io.Copy(twr, f); err != nil {
		fmt.Fprintf(os.Stderr,
			"error:[%v] writing file-%s to archive-%s",
			err, walInfo.Name(), archiveFile)
	}
	return f.Close()
}

func (w *Wal) getWalFileInfo() (fs.FileInfo, error) {
	f, err := os.Open(w.fileName)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr,
				"error:[%v] closing wal file-%s\n", err, w.WalFile())
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

// archiveFileName generates file name to be used for archive
// with timestamp numeric in unix micro format
// ex - archive-17188751512222238.tar
func archiveFileName() string {
	now := time.Now().UnixMicro()
	fname := fmt.Sprintf("archive-%d.tar", now)
	local := ".local"
	home := os.Getenv("HOME")
	fileName := filepath.Join(home, local, fname)
	return fileName
}
