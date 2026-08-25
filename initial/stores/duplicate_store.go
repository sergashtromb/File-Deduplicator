// хранилище для работы с общем массивом найденных дубликатов файлов

package stores

import (
	"bytes"
	"encoding/json"
	"file_deduplicator/domain"
	"log/slog"
	"slices"
	"sync"
)


type DeduplicateStore struct {
	rm sync.RWMutex
	dedupStore []*domain.Duplicate
}


func NewDeduplicateStore() *DeduplicateStore {
	return &DeduplicateStore{
		dedupStore: make([]*domain.Duplicate, 0),
	}
}


func (ds *DeduplicateStore) GetByHashAndSize(hash []byte, size int64) *domain.Duplicate {

	ds.rm.RLock()
	defer ds.rm.RUnlock()

	for _, val := range ds.dedupStore {

		if bytes.Equal(val.Hash, hash) && val.Size == size {
			return val
		}

	}

	return nil

}

func (ds *DeduplicateStore) Add(hash []byte, size int64, ff []*domain.FoundFile) *domain.Duplicate {

	ds.rm.Lock()
	defer ds.rm.Unlock()

	dup := domain.Duplicate {
		Hash: hash,
		Size: size,
		Files: ff,
	}

	ds.dedupStore = append(ds.dedupStore, &dup)

	return &dup

}

func (ds *DeduplicateStore) AddFoundFile(dup *domain.Duplicate, ff *domain.FoundFile) {

	ds.rm.Lock()
	defer ds.rm.Unlock()

	val := *dup
	have := slices.Contains(val.Files, ff)

	if !have {
		val.Files = append(val.Files, ff)
	}

}

func (ds *DeduplicateStore) GetJSONData() string {

	json_data, err := json.MarshalIndent(ds.dedupStore, "", "	")
	if err != nil {
		slog.Error("Errors Unmarshal json", "err", err)
	}

	res_string := string(json_data)

	return res_string

}