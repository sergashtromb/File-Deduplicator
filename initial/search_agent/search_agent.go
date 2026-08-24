// пакет для поиска файлов по наименованю

package search_agent

import (
	"context"
	"file_deduplicator/domain"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

type SearchAgent struct {
	WaitG sync.WaitGroup
}

func New() *SearchAgent {
	return &SearchAgent{
		WaitG: sync.WaitGroup{},
	}
}

func (sa *SearchAgent)FinderDirectories(ctx context.Context, params *domain.ParamsWorker, chPaths chan string) {

	sa.WaitG.Add(1)

	go func() {

		defer sa.WaitG.Done()
		defer close(chPaths)

		if params.UsRecurSeach {

			err := filepath.WalkDir(params.Path, func(path string, d fs.DirEntry, err error) error {

				if err != nil {
					slog.Error("Ошибка рекурсивного обхода директории")
					return err
				}

				if d.IsDir() {					
					chPaths <- path
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

				chPaths <- filepath.Join(params.Path, val.Name())

			}

		}

		select {
		case <- ctx.Done():
			slog.Debug("Завершение работы поиска папок")
			return
		}

	}()

}