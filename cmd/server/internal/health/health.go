package health

import (
	"go.uber.org/fx"
)

const (
	GroupTag = `group:"health"`
)

// Out annotates the return value of the constructor function with GroupTag.
func Out(fn any) any {
	return fx.Annotate(
		fn,
		fx.ResultTags(GroupTag),
	)
}

// In0 annotates the first argument of the function with GroupTag.
func In0(fn any) any {
	return fx.Annotate(
		fn,
		fx.ParamTags(GroupTag),
	)
}
