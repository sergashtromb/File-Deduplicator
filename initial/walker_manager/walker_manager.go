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
}

func New(limit int16, sa domain.SearchAgent, qg chan string) *WalkerManager {
	return &WalkerManager{
		QuantGr: limit,
		QueueGr: qg,
		SearchAgent: sa,
	}
}

func (wm *WalkerManager) StartGorutinesFindDuble(ctx context.Context, wg *sync.WaitGroup) {

	for i := 0; int16(i) < wm.QuantGr; i++ {

		wg.Add(1)

		go func() {
			// при завершении говорим о том что горутина выполнилась
			defer wg.Done()

			// бесконечный цикл чтения канала
			for {
				select {
				// считываем значения с канала
				case path, ok := <-wm.QueueGr:
					
					if ok == false {
						fmt.Println("DEAD")
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

