package url

import (
	"github.com/ervand7/urlshortener/internal/app/controllers/filetable"
	"sync"
)

type FileStorage struct {
	FileTable filetable.FileTable
	Mutex     sync.Mutex
}

func (f *FileStorage) Get(short string) (origin string, err error) {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	origin, err = f.FileTable.Get(short)
	if err != nil {
		return "", err
	}
	return origin, nil
}

func (f *FileStorage) Set(short, origin string) error {
	f.Mutex.Lock()
	defer f.Mutex.Unlock()
	if err := f.FileTable.Set(short, origin); err != nil {
		return err
	}
	return nil
}
