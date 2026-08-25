// файл нужен для настройки и инициализации механизма логирования

package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)


type DubleHandler struct {
	level 	slog.Leveler
	out 	io.Writer
}


func New(out io.Writer, level slog.Level) *DubleHandler {
	return &DubleHandler{out: out, level: level}

}

// проверяет выводить ли сообщение в логи, уровень задается в файлике конфигурации
func (h *DubleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *DubleHandler) Handle(ctx context.Context, r slog.Record) error {

	buf := make([]byte, 0, 1024)

	buf = append(buf, byte('['))
	buf = append(buf, []byte(LevelToString(r.Level))...)
	buf = append(buf, byte(']'))
	buf = append(buf, byte(' '))

	buf = append(buf, []byte(r.Time.Format("2006-01-02 15:04:05"))...)
	buf = append(buf, byte(' '))
	buf = append(buf, []byte(r.Message)...)

	r.Attrs(func(a slog.Attr) bool {
		buf = fmt.Appendf(buf, " %s=%v ", a.Key, a.Value)
		return true
	})

	buf = append(buf, byte('\n'))
	
	_, err := h.out.Write(buf)

	return err

}

func (h *DubleHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *DubleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// не поддерживаем атрибуты
	return h
}

func LevelToString(level slog.Level) string {

	a := "LVL"

	switch level {
	case slog.LevelInfo:
		a = "INFO"
	case slog.LevelDebug:
		a = "DEBUG"
	case slog.LevelWarn:
		a = "WARN"
	case slog.LevelError:
		a = "ERROR"
	default:
		a = "????"
	}
	
	return a

}

func StringToLevel(lvlstring *string) slog.Level {

	lvlstr := strings.ToLower(strings.TrimSpace(strings.Clone(*lvlstring)))

	var level slog.Level

	switch lvlstr {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}

	return level

}

func Init(path, lvlstring *string) (*os.File, error) {

	err := os.Mkdir(*path, 0755)
	if err != nil{}

	fullpath := filepath.Join(*path, "file_dd.log")

	file, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return  nil, err
	}

	level := StringToLevel(lvlstring)

	mlWriter := io.MultiWriter(os.Stdout, file)
	dubHandler := New(mlWriter, level)

	logger := slog.New(dubHandler)
	slog.SetDefault(logger)

	return file, nil

}