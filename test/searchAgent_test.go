package search_agent

import (
	//"file_deduplicator/domain"
	//"context"
	"file_deduplicator/domain"
	"file_deduplicator/initial/search_agent"
	"file_deduplicator/initial/stores"
	"path/filepath"
	"sync"
	"testing"
)

func BenchmarkDirectoryFinders(b *testing.B) {
	
	fileStore := stores.NewFileStore()
	dedupStore := stores.NewDeduplicateStore()

	searchAgent := search_agent.New(fileStore, dedupStore)

	params := domain.ParamsWorker {
		Path: filepath.Dir("Project"),
		UsRecurSeach: true,
		IsSaved: false,
	}

	var wg sync.WaitGroup

	b.ResetTimer()
	for i := 0; i < b.N; i++ {

		searchAgent.QueueGr = make(chan string)
		wg.Go(func() {
			select {
			case _, ok := <- searchAgent.QueueGr:
				if !ok {
					return
				}
			}
		})
		searchAgent.FinderDirectories(b.Context(), &params, &wg)	
		
		wg.Wait()
	}

}

func BenchmarkFindDubleFiles(b *testing.B) {

	fileStore := stores.NewFileStore()
	dedupStore := stores.NewDeduplicateStore()

	searchAgent := search_agent.New(fileStore, dedupStore)

	params := domain.ParamsWorker {
		Path: filepath.Dir("Project"),
		UsRecurSeach: true,
		IsSaved: false,
	}

	var wg sync.WaitGroup

	arr := make([]string, 0)

	wg.Go(func() {
		select {
		case val, ok := <- searchAgent.QueueGr:
			if !ok {
				return
			}
			arr = append(arr, val)
		}
	})

	searchAgent.FinderDirectories(b.Context(), &params, &wg)

	wg.Wait()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, arr_val := range arr {
			searchAgent.FinderFiles(arr_val)
		}
	}
	
}