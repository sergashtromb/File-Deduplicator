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
	QueueGr chan string
	fileStore domain.FileStore
}

func New(fs domain.FileStore) *SearchAgent {
	return &SearchAgent{
		QueueGr: make(chan string),
		fileStore: fs,
	}
}

func (sa *SearchAgent)FinderDirectories(ctx context.Context, params *domain.ParamsWorker, wg *sync.WaitGroup) {

	wg.Add(1)

	go func() {
		// TODO добавить учет контекста
		defer wg.Done()
		defer close(sa.QueueGr)

		if params.UsRecurSeach {

			err := filepath.WalkDir(params.Path, func(path string, d fs.DirEntry, err error) error {

				if err != nil {
					slog.Error("Ошибка рекурсивного обхода директории")
					return err
				}

				if d.IsDir() {					
					sa.QueueGr <- path
				}

				return nil

			})

			if err != nil {
				slog.Error("Ошибка чтения директории")
				return
			}

		} else {

			dirs, err := os.ReadDir(params.Path)
			if err != nil {
				slog.Error("Ошибка чтения директории")
				return
			}

			for _, val := range dirs {

				if val.IsDir() {
					sa.QueueGr <- filepath.Join(params.Path, val.Name())
				}
				
			}

		}

	}()
}

func (sa *SearchAgent) FinderFiles(path string) {

	vals, err := os.ReadDir(path)
	if err != nil {
		slog.Error("Ошибка считывания файлов из директории")
		return
	}

	for _, val := range vals {

		if val.IsDir() {
			continue
		}

		name := val.Name()
		f_info, err := val.Info()
		if err != nil {
			slog.Error("Ошибка считывания Info() файла")
			continue
		}

		f_path := filepath.Join(path, name)
		f_size := f_info.Size()

		hasher := sha256.New()

		file, err := os.Open(f_path)
		if err != nil {
			slog.Error("Ошибка открытия файла")
			continue
		}

		if _, err := io.Copy(hasher, file); err != nil {
			slog.Error("Ошибка вычисления хеша")
			file.Close()
			continue	
		}
		file.Close()

		hashsum := hasher.Sum(nil)

		ff := sa.fileStore.GetBySizeAndHash(f_size, hashsum)
		if ff != nil {
			slog.Debug("Найден дубль!")
		} else {
			sa.fileStore.Add(name, f_path, hashsum, f_size)
		}
		

	}

}