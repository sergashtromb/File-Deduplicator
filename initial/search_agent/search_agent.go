// пакет для поиска файлов по наименованю

package search_agent

import (
	"context"
	"crypto/sha256"
	"file_deduplicator/domain"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type SearchAgent struct {
	QueueGr 	chan string
	fileStore 	domain.FileStore
	dedupStore 	domain.DeduplicateStore
}

func New(fs domain.FileStore, ds domain.DeduplicateStore) *SearchAgent {
	return &SearchAgent{
		QueueGr: make(chan string),
		fileStore: fs,
		dedupStore: ds,
	}
}

func (sa *SearchAgent)FinderDirectories(ctx context.Context, params *domain.ParamsWorker, wg *sync.WaitGroup) {

	wg.Go(func() {

		defer close(sa.QueueGr)

		if params.UsRecurSeach {

			err := filepath.WalkDir(params.Path, func(path string, d fs.DirEntry, err error) error {

				if err != nil {
					slog.Error("Error recursivy directory traversal")
					return err
				}

				if d.IsDir() {					
					sa.QueueGr <- path
				}

				return nil

			})

			if err != nil {
				slog.Error("Error directory read", "err", err)
				return
			}

		} else {

			dirs, err := os.ReadDir(params.Path)
			if err != nil {
				slog.Error("Error directory read", "err", err)
				return
			}

			for _, val := range dirs {

				if val.IsDir() {
					sa.QueueGr <- filepath.Join(params.Path, val.Name())
				}
				
			}

		}

	})
}

func (sa *SearchAgent) FinderFiles(path string) {

	vals, err := os.ReadDir(path)
	if err != nil {
		slog.Error("Error file read from directory", "err", err)
		return
	}

	for _, val := range vals {

		if val.IsDir() {
			continue
		}

		name := val.Name()
		f_info, err := val.Info()
		if err != nil {
			slog.Error("Error use Info()", "err", err)
			continue
		}

		f_path := filepath.Join(path, name)
		f_size := f_info.Size()

		hasher := sha256.New()

		file, err := os.OpenFile(f_path, os.O_RDONLY, 0444)
		if err != nil {
			slog.Error("Error file open", "err", err)
			continue
		}

		if _, err := io.Copy(hasher, file); err != nil {
			slog.Error("Error hash calculation")
			file.Close()
			continue	
		}
		file.Close()

		hashsum := hasher.Sum(nil)

		new_ff := sa.fileStore.Add(name, f_path, hashsum, f_size)
		ff := sa.fileStore.GetBySizeAndHash(new_ff.Size, new_ff.Hash)

		if ff != nil && new_ff != ff {
			dup := sa.dedupStore.GetByHashAndSize(hashsum, f_size)

			if dup == nil {
				arr_ff := []*domain.FoundFile{new_ff, ff}
				_ = sa.dedupStore.Add(hashsum, f_size, arr_ff)

			} else {
				sa.dedupStore.AddFoundFile(dup, new_ff)

			}

		}
		

	}

}