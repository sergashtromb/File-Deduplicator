// пакет для отчета о найденных файлах выводит результат в консоль или отдельный json файл
package report

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Report(jsonString string, isSaved bool, reportSavePath string) {

	if !isSaved {
		fmt.Println(jsonString)
	} else {
		time := time.Now()
		compFileName := []string{"report_", time.Format("2006_01_02_15_04_05"), ".json"}

		fileName := strings.Join(compFileName, "")

		err := os.WriteFile(filepath.Join(reportSavePath, fileName), []byte(jsonString), 0755)
		if err != nil {
			slog.Error("Error write report", "err", err)
		} 

	}

}