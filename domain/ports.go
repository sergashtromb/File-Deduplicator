// модуль предназначен для описания интерфейсов

package domain

import (
	"context"
)

// описание менеджера
//	StartGorutines - считывает канал очереди и запускает горутину если лимиты позволяют, завершает работу если канал закроется
//	AddInQueueGr - добавляет запись в канал очереди.
type WalkerManager interface {
	StartGorutinesFindDuble(ctx context.Context)
	AddInQueueGr(path string)
}

type SearchAgent interface {
	//LoadFilesData(path string)
	FinderDirectories(ctx context.Context, params *ParamsWorker, chPaths chan string)
}