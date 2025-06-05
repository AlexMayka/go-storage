package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
}

type SLogger struct {
	logger *slog.Logger
}

func (l SLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l SLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l SLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l SLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func getSlogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getIoWriter(saveToFile bool) (io.Writer, error) {
	var w io.Writer = os.Stdout

	if saveToFile {
		if err := os.MkdirAll("log", os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %v", err)
		}

		file, err := os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %v", err)
		}
		w = file
	}
	return w, nil
}

func InitLog(saveToFile bool, level string) (Logger, error) {
	w, err := getIoWriter(saveToFile)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	levelSlog := getSlogLevel(level)
	return SLogger{logger: slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: levelSlog,
	}))}, nil
}

func WithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, "logger", l)
}

func FromContext(ctx context.Context) Logger {
	if logging, ok := ctx.Value("logger").(Logger); ok {
		return logging
	}

	return SLogger{logger: slog.Default()}
}
