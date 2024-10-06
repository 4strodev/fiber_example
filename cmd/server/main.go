package main

import (
	"log/slog"

	"github.com/4strodev/go_monitoring_example/pkg/features/auth"
	"github.com/4strodev/go_monitoring_example/pkg/shared"
	"github.com/4strodev/go_monitoring_example/pkg/shared/events"
	scaffold "github.com/4strodev/scaffold/pkg/core/app"
	wiring "github.com/4strodev/wiring/pkg"
)

func main() {
	container := wiring.New()
	container.Transient(auth.NewAuthService)
	container.Singleton(shared.NewDBClient)
	container.Singleton(events.NewEventBus)
	container.Singleton(shared.NewLogger)
	app := scaffold.NewApp(container)

	app.AddAdapter(&shared.FiberAdapter{})
	app.AddComponent(&auth.AuthController{})

	var logger *slog.Logger
	err := container.Resolve(&logger)
	if err != nil {
		panic(err)
	}

	// TODO handle shutdown should be a method for app
	err = app.Start()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
