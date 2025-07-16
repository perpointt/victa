package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"
	"time"
)

/* ---------- уровни ---------- */

type Level uint8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelOff
)

var levelNames = [...]string{"DEBUG", "INFO", "WARN", "ERROR"}

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)

	SetLevel(lvl Level)
	Level() Level
}

type stdLogger struct {
	out   *log.Logger
	level atomic.Uint32
}

// New возвращает потокобезопасный Logger.
// writer — куда писать (nil = os.Stderr), lvl — минимальный уровень вывода.
func New(writer io.Writer, lvl Level) Logger {
	if writer == nil {
		writer = os.Stderr
	}
	l := &stdLogger{
		out: log.New(writer, "", 0), // флаги не нужны, дату вставляем сами
	}
	l.level.Store(uint32(lvl))
	return l
}

func (l *stdLogger) Debug(f string, a ...any) { l.log(LevelDebug, f, a...) }
func (l *stdLogger) Info(f string, a ...any)  { l.log(LevelInfo, f, a...) }
func (l *stdLogger) Warn(f string, a ...any)  { l.log(LevelWarn, f, a...) }
func (l *stdLogger) Error(f string, a ...any) { l.log(LevelError, f, a...) }

func (l *stdLogger) SetLevel(lvl Level) { l.level.Store(uint32(lvl)) }
func (l *stdLogger) Level() Level       { return Level(l.level.Load()) }

func (l *stdLogger) log(lvl Level, format string, args ...any) {
	if lvl < Level(l.level.Load()) {
		return
	}

	_, file, line, _ := runtime.Caller(3) // 3 — Debug/Info → log → runtime
	prefix := fmt.Sprintf("%s %s:%d [%s] ",
		time.Now().UTC().Format(time.RFC3339),
		filepath.Base(file), line,
		levelNames[lvl],
	)

	l.out.Printf(prefix+format, args...)
}
