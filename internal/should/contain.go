package should

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/samber/lo"
)

var (
	_ validation.Rule = (*ContainRule[string])(nil)
)

// Contain returns a validation rule to check if the slice value contains the required elements. The rule ignores
// any input value that is not a slice whose element type is the required elements type.
func Contain[T comparable](elements ...T) *ContainRule[T] {
	return &ContainRule[T]{
		required: elements,
		message:  "should contain required elements",
	}
}

type ContainRule[T comparable] struct {
	required []T
	message  string
}

func (r *ContainRule[T]) Error(message string) *ContainRule[T] {
	r.message = message
	return r
}

func (r *ContainRule[T]) Validate(value interface{}) error {
	slice, ok := value.([]T)
	if !ok {
		return nil
	}

	if !lo.Every(slice, r.required) {
		return errors.New(r.message)
	}

	return nil
}
