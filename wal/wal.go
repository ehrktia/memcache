package wal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ehrktia/memcache/datastructure"
)

// receive data and persist in file
// keep track of pointer up to where the data is being loaded in to cache

type Wal struct {
	fileName
}

type fileName struct {
	defaultWalDir   string
	fileName        string
	archiveFileName string
}

// createFile checks if file exist
// and create new file when not found
func CreateFile(filePath string) error {
	if !exists(filePath) {
		_, err := os.Create(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Wal) WalFileName() string {
	return w.fileName.fileName
}

func (w *Wal) updWalFile(file string) {
	walFile := filepath.Join(w.defaultWalDir, file)
	w.fileName.fileName = walFile
}

func genWalDir() string {
	home := os.Getenv("HOME")
	local := os.Getenv("LOCAL")
	var walFileDir string
	if local == "" {
		local = ".local"
		walFileDir = filepath.Join(home, local)
		return walFileDir
	} else {
		walFileDir = filepath.Join(home, local)
		return walFileDir
	}
}

func genArchiveDir() string {
	archiveDir := os.Getenv("ARCHIVE_DIR")
	if archiveDir == "" {
		return genWalDir()
	}
	return archiveDir
}

// creates file if one does not exist
func NewWal() *Wal {
	var defaultDir string
	var err error
	if os.Getenv("CI") != "" {
		defaultDir, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "get dir err:%v\n", err)
			os.Exit(1)
		}

	} else {
		defaultDir = genWalDir()
	}
	defaultMaxFileSize = getDefaultFileSize()
	now := time.Now().UnixMicro()
	archiveDir := genArchiveDir()
	archiveName := fmt.Sprintf("archive-%d.tar", now)
	walFileName := filepath.Join(defaultDir, fmt.Sprintf("wal-%d.txt", now))
	archiveFileName := filepath.Join(archiveDir, archiveName)
	return &Wal{
		fileName{
			defaultWalDir:   defaultDir,
			fileName:        walFileName,
			archiveFileName: archiveFileName,
		},
	}
}

func exists(fileName string) bool {
	f, err := os.Open(fileName)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	if err := f.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "error:[%v] closing-%q file\n", err, fileName)
		os.Exit(1)
	}
	return true
}

var ErrEmptyFile = errors.New("empty file")
var ErrRead = errors.New("no data to read")

// UpdCache writes data into wal for persisence
// and triggers cache update, adds data in to cache
func UpdCache(w *Wal, data []byte) error {
	if _, err := Write(w.fileName.fileName, data); err != nil {
		return err
	}
	// trigger cache update
	d := new(datastructure.Data)
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}
	v, isLoaded := datastructure.Add(d.Key, d.Value)
	if isLoaded {
		fmt.Fprintf(os.Stdout, "[INFO]: cache value already found, updated with:%v\n", v)
	}
	return nil
}
