package main

import (
	"context"
	"file_deduplicator/config"
	"file_deduplicator/domain"
	"file_deduplicator/initial/search_agent"
	"file_deduplicator/initial/walker_manager"
	"file_deduplicator/logger"
	"fmt"
	"log/slog"
	"sync"

	//"log/slog"
	"os"
)

func main() {
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	params := parseParamsWorker(os.Args[1:])

	_, err := os.Stat(params.Path)

	if err != nil {
		slog.Error("Такого пути не существует")
	}

	cnf := config.Load("config.yaml")

	file, err := logger.Init(&cnf.LogSettings.Directory, &cnf.LogSettings.Level)
	if err != nil {
		fmt.Println("Ошибка формирования логов")
	}
	defer file.Close()

	var wg sync.WaitGroup

	searchAgent := search_agent.New()
	walkerManager := walker_manager.New(10, searchAgent, searchAgent.QueueGr)

	walkerManager.StartGorutinesFindDuble(ctx, &wg)
	searchAgent.FinderDirectories(ctx, params, &wg)
	
	wg.Wait()


}

func parseParamsWorker(args []string) *domain.ParamsWorker {

	defParams := domain.ParamsWorker{
		IsSaved: false,
		UsRecurSeach: false,
		Path: "",
	}

	for _, val := range args {

		switch val {
		case "-s":
			defParams.IsSaved = true
		case "-r":
			defParams.UsRecurSeach = true
		default:
			defParams.Path = val

		}

	}

	return &defParams

}