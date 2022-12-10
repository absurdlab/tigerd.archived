package spec

import (
	"encoding/json"
	"fmt"
)

const (
	AppTypeWeb ApplicationType = 1 << iota
	AppTypeSinglePage
	AppTypeNative
	AppTypeMachine

	appTypeWeb     = "web"
	appTypeSpa     = "spa"
	appTypeNative  = "native"
	appTypeMachine = "machine"
)

// ApplicationType represents application_type parameter in OpenID Connect 1.0 dynamic registration.
type ApplicationType uint8

func (t ApplicationType) String() string {
	switch t {
	case AppTypeWeb:
		return appTypeWeb
	case AppTypeSinglePage:
		return appTypeSpa
	case AppTypeNative:
		return appTypeNative
	case AppTypeMachine:
		return appTypeMachine
	default:
		return ""
	}
}

func (t ApplicationType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *ApplicationType) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case appTypeWeb:
		*t = AppTypeWeb
	case appTypeSpa:
		*t = AppTypeSinglePage
	case appTypeNative:
		*t = AppTypeNative
	case appTypeMachine:
		*t = AppTypeMachine
	default:
		return fmt.Errorf("invalid value for spec.ApplicationType [%s]", value)
	}

	return nil
}
