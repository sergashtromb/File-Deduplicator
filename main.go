package main

import (
	"file_deduplicator/logger"
	"fmt"
	"log/slog"
)

func main() {
	
	file, err := logger.Init()
	if err != nil {
		fmt.Println("Ошибка формирования логов")
	}
	defer file.Close()

	slog.Debug("logger инициализирован")
	slog.Info("Start searching...")

	slog.Debug("Завершение поиска")


}