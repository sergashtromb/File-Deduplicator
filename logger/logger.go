// файл нужен для настройки и инициализации механизма логирования

package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)


type DubleHandler struct {
	level slog.Leveler
	out io.Writer
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

func Init() (*os.File, error) {

	err := os.Mkdir("logs", 0755)
	if err != nil{}

	file, err := os.OpenFile("logs/file_dd.log", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return  nil, err
	}

	mlWriter := io.MultiWriter(os.Stdout, file)
	dubHandler := New(mlWriter, slog.LevelDebug)

	logger := slog.New(dubHandler)
	slog.SetDefault(logger)

	return file, nil

}