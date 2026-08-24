// пакет для контроля запуска горутин поиска, а так же их счетчика

package walker_manager

import (
	"context"
	"file_deduplicator/domain"
	"fmt"
	"log/slog"
	"sync"
)

// описание структуры менеджера контроля количества горутин
//	quant_gr - количество горутин
//	queue_gr - очередь директорий

type WalkerManager struct {
	QuantGr 	int16
	QueueGr 	chan string
	SearchAgent domain.SearchAgent
	WaitG 		sync.WaitGroup
}

func New(limit int16, sa domain.SearchAgent) *WalkerManager {
	return &WalkerManager{
		QuantGr: limit,
		QueueGr: make(chan string, 1000),
		WaitG: sync.WaitGroup{},
		SearchAgent: sa,
	}
}

func (wm *WalkerManager) StartGorutinesFindDuble(ctx context.Context) {

	for i := 0; int16(i) < wm.QuantGr; i++ {

		wm.WaitG.Add(1)

		go func() {
			// при завершении говорим о том что горутина выполнилась
			defer wm.WaitG.Done()

			// бесконечный цикл чтения канала
			for {
				select {
				// считываем значения с канала
				case path, ok := <-wm.QueueGr:

					if !ok {
						return
					}
					fmt.Println(path)
				// завершаем при остановке
				case <-ctx.Done():
					
					slog.Debug("Завершаем поиск!")
					return

				}
			}

			//wm.SearchAgent.LoadFilesData()

		}()
	}

}

