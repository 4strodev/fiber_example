package shared

import (
	"io"
	"log/slog"
	"os"
	"path"
)

const log_file = "./tmp/log/logs.txt"

func NewLogger() (*slog.Logger, error) {
	directory := path.Dir(log_file)
	err := os.MkdirAll(directory, os.ModePerm|os.ModeDir)
	if err != nil {
		return nil, err
	}
	logFile, err := os.OpenFile(log_file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	attributes := []slog.Attr{
		slog.Any("service_name", "fiber_example"),
	}

	output := io.MultiWriter(logFile, os.Stdout)

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}).WithAttrs(attributes)
	return slog.New(handler), nil
}
