package handler

import (
	"github.com/absurdlab/tigerd/internal/jose"
	"github.com/absurdlab/tigerd/internal/wellknown"
	"github.com/labstack/echo/v4"
	"net/http"
)

func NewWellKnownHandler(discovery *wellknown.Discovery, jwks *jose.JSONWebKeySet) Interface {
	return &wellKnownHandler{
		discovery: discovery,
		jwks:      jwks,
	}
}

type wellKnownHandler struct {
	discovery *wellknown.Discovery
	jwks      *jose.JSONWebKeySet
}

func (h *wellKnownHandler) Mount(e *echo.Echo) error {
	e.GET("/.well-known/openid-configuration", h.getDiscovery)
	e.GET("/.well-known/jwks.json", h.getJSONWebKeySet)

	return nil
}

func (h *wellKnownHandler) getDiscovery(e echo.Context) error {
	return e.JSON(http.StatusOK, h.discovery)
}

func (h *wellKnownHandler) getJSONWebKeySet(e echo.Context) error {
	e.Response().Header().Set(echo.HeaderContentType, "application/jwk-set+json")
	return e.JSON(http.StatusOK, h.jwks.Public())
}
