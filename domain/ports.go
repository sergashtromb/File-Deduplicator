// модуль предназначен для описания интерфейсов

package domain

import (
	"context"
	"sync"
)

// описание менеджера
//	StartGorutines - считывает канал очереди и запускает горутину если лимиты позволяют, завершает работу если канал закроется
//	AddInQueueGr - добавляет запись в канал очереди.
type WalkerManager interface {
	StartGorutinesFindDuble(ctx context.Context, wg *sync.WaitGroup)
	AddInQueueGr(path string)
}

type SearchAgent interface {
	//LoadFilesData(path string)
	FinderDirectories(ctx context.Context, params *ParamsWorker, wg *sync.WaitGroup)
	FinderFiles(path string)
}

type FileStore interface {
	GetBySizeAndHash(size int64, hash []byte) *FoundFile
	Add(name, path string, hash []byte, size int64)
}