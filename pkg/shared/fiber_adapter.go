package shared

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/4strodev/wiring/pkg"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type FiberAdapter struct {
	app    *fiber.App
	logger *slog.Logger
}

// Init implements adapters.Adapter.
func (f *FiberAdapter) Init(container pkg.Container) error {
	err := container.Resolve(&f.logger)
	if err != nil {
		return nil
	}
	f.app = fiber.New(fiber.Config{
		ErrorHandler: f.errorHandler,
		ReadTimeout:  time.Second * 30,
	})

	f.app.Use(recover.New())
	f.app.Use(func(ctx fiber.Ctx) error {
		var logLevel slog.Level = slog.LevelInfo
		response := ctx.Response()
		errChain := ctx.Next()
		if errChain != nil {
			logLevel = slog.LevelError

			if err := ctx.App().ErrorHandler(ctx, errChain); err != nil {
				ctx.SendStatus(http.StatusInternalServerError)
			}
		}

		statusCode := response.StatusCode()
		method := ctx.Method()
		path := ctx.Path()

		f.logger.Log(context.Background(), logLevel, "{method} {path} {statusCode}", "method", method, "path", path, "statusCode", statusCode)

		return nil
	})

	f.app.Hooks().OnListen(func(data fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		f.logger.Info("listening on port 8080")
		return nil
	})
	container.Singleton(func() fiber.Router {
		return f.app
	})
	return nil
}

// Start implements adapters.Adapter.
func (f *FiberAdapter) Start() error {
	return f.app.Listen(":8080", fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}

// Stop implements adapters.Adapter.
func (f *FiberAdapter) Stop() error {
	return f.app.Shutdown()
}

func (f *FiberAdapter) errorHandler(ctx fiber.Ctx, err error) error {
	f.logger.Error(err.Error())
	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"msg": err.Error(),
	})
}
