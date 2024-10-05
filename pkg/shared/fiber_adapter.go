package shared

import (
	"log/slog"
	"net/http"

	"github.com/4strodev/wiring/pkg"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

type FiberAdapter struct {
	app *fiber.App
}

// Init implements adapters.Adapter.
func (f *FiberAdapter) Init(container pkg.Container) error {
	var log *slog.Logger
	err := container.Resolve(&log)
	if err != nil {
		return nil
	}
	f.app = fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	f.app.Use(recover.New())
	f.app.Use(func(ctx fiber.Ctx) error {
		response := ctx.Response()
		errChain := ctx.Next()
		if errChain != nil {
			if err := ctx.App().ErrorHandler(ctx, errChain); err != nil {
				ctx.SendStatus(http.StatusInternalServerError)
			}
		}

		statusCode := response.StatusCode()
		method := ctx.Method()
		path := ctx.Path()

		log.Info("{method} {path} {statusCode}", "method", method, "path", path, "statusCode", statusCode)

		return errChain
	})

	f.app.Hooks().OnListen(func(data fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		log.Info("listening on port 3000")
		return nil
	})
	container.Singleton(func() fiber.Router {
		return f.app
	})
	return nil
}

// Start implements adapters.Adapter.
func (f *FiberAdapter) Start() error {
	return f.app.Listen(":3000", fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}

// Stop implements adapters.Adapter.
func (f *FiberAdapter) Stop() error {
	return f.app.Shutdown()
}

func errorHandler(ctx fiber.Ctx, err error) error {
	return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"msg": err.Error(),
	})
}
