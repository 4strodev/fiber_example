package main

import (
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
	container.Singleton(func() *slog.Logger {
		attributes := []slog.Attr{
			slog.Any("source", "fiber_example"),
		}
		handler := slog.NewJSONHandler(os.Stdout, nil).WithAttrs(attributes)
		return slog.New(handler)
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
