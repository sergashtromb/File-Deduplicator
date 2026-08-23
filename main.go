package main

import (
	"file_deduplicator/logger"
	"file_deduplicator/config"
	"fmt"
	"log/slog"
)

func main() {
	
	cnf := config.Load("config.yaml")

	file, err := logger.Init(&cnf.LogSettings.Directory, &cnf.LogSettings.Level)
	if err != nil {
		fmt.Println("Ошибка формирования логов")
	}
	defer file.Close()

	slog.Debug("logger инициализирован")
	slog.Info("Start searching...")

	slog.Debug("Завершение поиска")


}