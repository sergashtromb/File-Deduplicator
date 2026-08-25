// пакет для контроля запуска горутин поиска, а так же их счетчика

package walker_manager

import (
	"context"
	"file_deduplicator/domain"
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

		wg.Go(func() {
		
			// бесконечный цикл чтения канала
			for {

				select {
				// считываем значения с канала
				case path, ok := <-wm.QueueGr:
					
					if !ok {
						return
					}
					
					wm.SearchAgent.FinderFiles(path)
				
				case <-ctx.Done():
					
					slog.Debug("Finish search!")
					return

				}
			}

		})
	}

}

