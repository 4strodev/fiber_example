package auth

import (
	"log/slog"

	"github.com/4strodev/wiring/pkg"
	"github.com/gofiber/fiber/v3"
)

type AuthController struct {
	Service AuthService
	Logger  *slog.Logger
	Router  fiber.Router
}

// Init implements components.Component.
func (a *AuthController) Init(container pkg.Container) error {
	err := container.Fill(a)
	if err != nil {
		return err
	}

	// Setting routes
	group := a.Router.Group("/auth")
	group.Post("/login", a.Login)
	group.Post("/register", a.Register)

	return nil
}

func (c *AuthController) Register(ctx fiber.Ctx) error {
	var requestBody RegisterRequest
	err := ctx.Bind().Body(&requestBody)
	if err != nil {
		return err
	}

	err = c.Service.Register(requestBody)
	if err != nil {
		return err
	}

	return ctx.JSON(fiber.Map{
		"msg": "user registered",
	})
}

func (c *AuthController) Login(ctx fiber.Ctx) error {
	requestBody := LoginRequest{}
	err := ctx.Bind().Body(&requestBody)
	if err != nil {
		return err
	}
	response, err := c.Service.Login(requestBody)
	if err != nil {
		return err
	}
	return ctx.JSON(fiber.Map{
		"data": response,
		"msg":  "everything okay",
	})
}
