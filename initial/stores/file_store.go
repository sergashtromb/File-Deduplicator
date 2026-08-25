// хранилище для работы с общим массивом всех уже найденных файлов

package stores

import (
	"bytes"
	"file_deduplicator/domain"
	"sync"
)

type FileStore struct {
	rm sync.RWMutex
	filesData []*domain.FoundFile
}

func NewFileStore() *FileStore {
	return &FileStore{
		filesData: make([]*domain.FoundFile, 0),
	}
}

func (fs *FileStore) GetBySizeAndHash(size int64, hash []byte) *domain.FoundFile {

	fs.rm.RLock()
	defer fs.rm.RUnlock()

	for _, dt := range fs.filesData {

		if bytes.Equal(dt.Hash, hash) && dt.Size == size {
			return dt
		}

	}

	return nil

}

func (fs *FileStore) Add(name, path string, hash []byte, size int64) *domain.FoundFile {

	fs.rm.Lock()
	defer fs.rm.Unlock()

	ff := domain.FoundFile {
		Name: name,
		Hash: hash,
		Path: path,
		Size: size,
	}

	fs.filesData = append(fs.filesData, &ff)
	
	return &ff

}