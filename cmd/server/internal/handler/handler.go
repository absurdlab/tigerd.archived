package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

const (
	GroupTag = `group:"handlers"`
)

// H is implemented by all api endpoints for handling http/grpc traffic.
type H interface {
	// Mount mounts the endpoints available under this implementation of H to the echo framework.
	Mount(e *echo.Echo) error
}

// Out returns an annotated function where the result of the function is tagged with GroupTag.
func Out(fn any) any {
	return fx.Annotate(
		fn,
		fx.As(new(H)),
		fx.ResultTags(GroupTag),
	)
}

// In0 returns an annotated function where the first parameter of the function is tagged with GroupTag.
func In0(fn any) any {
	return fx.Annotate(
		fn,
		fx.ParamTags(GroupTag),
	)
}
