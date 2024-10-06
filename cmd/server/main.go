package main

import (
	"io"
	"log/slog"
	"os"
	"os/signal"

	"github.com/4strodev/go_monitoring_example/pkg/features/auth"
	"github.com/4strodev/go_monitoring_example/pkg/shared"
	"github.com/4strodev/scaffold/pkg/core/app"
	wiring "github.com/4strodev/wiring/pkg"
)

func main() {
	container := wiring.New()
	container.Transient(auth.NewAuthService)
	container.Singleton(shared.NewDBClient)
	container.Singleton(func() (*slog.Logger, error) {
		err := os.MkdirAll("/tmp/log/fiber_example", os.ModePerm|os.ModeDir)
		if err != nil {
			return nil, err
		}
		logFile, err := os.OpenFile("/tmp/log/fiber_example/logs.txt", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
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
	})
	app := app.NewApp(container)

	app.AddAdapter(&shared.FiberAdapter{})
	app.AddComponent(&auth.AuthController{})

	var logger *slog.Logger
	err := container.Resolve(&logger)
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	errs, err := app.Start()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	for err := range errs {
		logger.Error(err.Error())
	}

	<-sigs
	err = app.Stop()
	if err != nil {
		logger.Error(err.Error())
	}
}
