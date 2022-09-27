package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

type Local struct {
	maxFileSize int    // max number of bytes for files
	basePath    string // path of storage
}

func NewLocal(basePath string, maxSize int) (*Local, error) {
	path, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &Local{basePath: path}, nil
}

func (l *Local) Save(path string, contents io.Reader) error {

	fp := l.fullPath(path)

	// gets the dir if not exists create one
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}

	// if file exists delete it
	_, err = os.Stat(fp)
	if err == nil {
		err = os.Remove(fp)
		if err != nil {
			return xerrors.Errorf("Unable to delete file: %w", err)
		}
	}

	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()

	// write the contents to the new file
	_, err = io.Copy(f, contents)
	if err != nil {
		return xerrors.Errorf("Unable to write to file: %w", err)
	}

	return nil
}

func (l *Local) Get(path string) (*os.File, error) {
	fp := l.fullPath(path)

	f, err := os.Open(fp)
	if err != nil {
		return nil, xerrors.Errorf("Unable to open the file: %w", err)
	}

	fmt.Println("[DEBUG] entered the Get local")
	return f, nil
}

// return absolute path
func (l *Local) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}
