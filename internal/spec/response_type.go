package spec

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	ResponseTypeCode ResponseType = 1 << iota
	ResponseTypeToken
	ResponseTypeIDToken

	responseTypeCode    = "code"
	responseTypeToken   = "token"
	responseTypeIDToken = "id_token"
)

// ResponseType represents the response_type parameter defined in OAuth 2.0 and OpenID Connect 1.0.
type ResponseType uint8

func (t ResponseType) UInt8() uint8 {
	return uint8(t)
}

func (t ResponseType) String() string {
	switch t {
	case ResponseTypeCode:
		return responseTypeCode
	case ResponseTypeToken:
		return responseTypeToken
	case ResponseTypeIDToken:
		return responseTypeIDToken
	default:
		return ""
	}
}

// NewResponseTypeSet creates a new ResponseTypeSet with given response types set in it.
func NewResponseTypeSet(responseType ...ResponseType) *ResponseTypeSet {
	set := ResponseTypeSet(0).Ref()
	for _, t := range responseType {
		set.Add(t)
	}
	return set
}

// ResponseTypeSet is a bitmask of zero or more ResponseType that models multi response type format. It is distinct
// from an array of ResponseType. Instead, it marshals into a space delimited string containing multiple response
// type values.
type ResponseTypeSet uint8

func (s ResponseTypeSet) UInt8() uint8 {
	return uint8(s)
}

func (s ResponseTypeSet) Ref() *ResponseTypeSet {
	return &s
}

func (s ResponseTypeSet) Contains(t ResponseType) bool {
	return s.UInt8()&t.UInt8() != 0
}

func (s *ResponseTypeSet) Add(t ResponseType) *ResponseTypeSet {
	*s = ResponseTypeSet(s.UInt8() | t.UInt8())
	return s
}

func (s ResponseTypeSet) MarshalJSON() ([]byte, error) {
	if s == 0 {
		return nil, nil
	}

	var values []string
	for _, t := range []ResponseType{ResponseTypeCode, ResponseTypeToken, ResponseTypeIDToken} {
		if s.Contains(t) {
			values = append(values, t.String())
		}
	}

	return json.Marshal(strings.Join(values, " "))
}

func (s *ResponseTypeSet) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return nil
	}

	value = strings.TrimSpace(value)
	if len(value) == 0 {
		return nil
	}

	for _, elem := range strings.Split(value, " ") {
		switch elem {
		case responseTypeCode:
			s.Add(ResponseTypeCode)
		case responseTypeToken:
			s.Add(ResponseTypeToken)
		case responseTypeIDToken:
			s.Add(ResponseTypeIDToken)
		default:
			return fmt.Errorf("invalid spec.ResponseType value [%s]", elem)
		}
	}

	return nil
}
