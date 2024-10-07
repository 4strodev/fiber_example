package main

import (
	"encoding/json"
	"log/slog"

	"github.com/4strodev/go_monitoring_example/pkg/features/auth"
	"github.com/4strodev/go_monitoring_example/pkg/shared"
	"github.com/4strodev/go_monitoring_example/pkg/shared/events"
	scaffold "github.com/4strodev/scaffold/pkg/core/app"
	wiring "github.com/4strodev/wiring/pkg"
)

func main() {
	var err error
	container := wiring.New()
	container.Transient(auth.NewAuthService)
	container.Singleton(shared.NewDBClient)
	container.Singleton(events.NewEventBus)
	container.Singleton(shared.NewLogger)

	app := scaffold.NewApp(container)

	app.AddAdapter(&shared.FiberAdapter{})
	app.AddComponent(&auth.AuthController{})

	var logger *slog.Logger
	err = container.Resolve(&logger)
	if err != nil {
		panic(err)
	}

	var eventBus events.EventBus
	err = container.Resolve(&eventBus)
	if err != nil {
		panic(err)
	}

	err = eventBus.Listen("*", func(event events.Event) error {
		json, err := json.Marshal(event)
		if err != nil {
			return err
		}

		logger.Info("new event", "event", string(json))
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = app.Start()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}
