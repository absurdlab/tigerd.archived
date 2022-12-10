package spec

import (
	"encoding/json"
	"fmt"
)

const (
	ResponseModeQuery ResponseMode = 1 << iota
	ResponseModeFragment

	responseModeQuery    = "query"
	responseModeFragment = "fragment"
)

// ResponseMode represents response_mode parameter in OpenID Connect 1.0.
type ResponseMode uint8

func (m ResponseMode) String() string {
	switch m {
	case ResponseModeQuery:
		return responseModeQuery
	case ResponseModeFragment:
		return responseModeFragment
	default:
		return ""
	}
}

func (m ResponseMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *ResponseMode) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case responseModeQuery:
		*m = ResponseModeQuery
	case responseModeFragment:
		*m = ResponseModeFragment
	default:
		return fmt.Errorf("invalid value for spec.ResponseMode [%s]", value)
	}

	return nil
}
