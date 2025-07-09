package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

// Logger — наш интерфейс для логирования.
type Logger interface {
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(format string, args ...interface{})

	InfoDepth(depth int, format string, args ...interface{})
	WarnDepth(depth int, format string, args ...interface{})
	ErrorDepth(depth int, format string, args ...interface{})
	DebugDepth(depth int, format string, args ...interface{})
}

// stdLogger — конкретная реализация на базе log.Logger.
type stdLogger struct {
	base *log.Logger
}

// New создает Logger, выводящий дату, время и файл:строку.
func New() Logger {
	// логгер по умолчанию пишет в stderr, флаги будут устанавливаться вручную.
	base := log.New(os.Stderr, "", 0)
	// устанавливаем формат: дата+время + краткий файл/строка
	base.SetFlags(log.LstdFlags) // yyyy/mm/dd hh:mm:ss
	return &stdLogger{
		base: base,
	}
}

// общий вывод
func (l *stdLogger) output(depth int, level, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	prefix := fmt.Sprintf("%s:%d [%s] ", filepath.Base(file), line, level)
	msg := fmt.Sprintf(format, args...)
	// 0 тут значит — не добавлять ещё раз дату/время (они уже в base.Flags)
	l.base.Output(0, prefix+msg)
}

func (l *stdLogger) Info(format string, args ...interface{})   { l.output(2, "INFO", format, args...) }
func (l *stdLogger) Errorf(format string, args ...interface{}) { l.output(2, "ERROR", format, args...) }
func (l *stdLogger) Warn(format string, args ...interface{})   { l.output(2, "WARN", format, args...) }
func (l *stdLogger) Debug(format string, args ...interface{})  { l.output(2, "DEBUG", format, args...) }

func (l *stdLogger) InfoDepth(depth int, format string, args ...interface{}) {
	l.output(depth, "INFO", format, args...)
}
func (l *stdLogger) WarnDepth(depth int, format string, args ...interface{}) {
	l.output(depth, "WARN", format, args...)
}
func (l *stdLogger) ErrorDepth(depth int, format string, args ...interface{}) {
	l.output(depth, "ERROR", format, args...)
}
func (l *stdLogger) DebugDepth(depth int, format string, args ...interface{}) {
	l.output(depth, "DEBUG", format, args...)
}

// itoa — быстрый int→string
func itoa(i int) string {
	return strconv.FormatInt(int64(i), 10)
}
