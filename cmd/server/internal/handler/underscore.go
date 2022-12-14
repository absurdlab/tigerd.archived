package handler

import (
	"github.com/hellofresh/health-go/v5"
	"github.com/labstack/echo/v4"
)

func NewUnderscoreHandler(health *health.Health) Interface {
	return &underscoreHandler{
		health: health,
	}
}

type underscoreHandler struct {
	health *health.Health
}

func (h *underscoreHandler) Mount(e *echo.Echo) error {
	e.GET("/_/health", echo.WrapHandler(h.health.Handler()))

	return nil
}
