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

func (t ResponseType) ToSet() ResponseTypeSet {
	return ResponseTypeSet(0).Add(t)
}

// ResponseTypeSet is a bitmask of zero or more ResponseType that models multi response type format. It is distinct
// from an array of ResponseType. Instead, it marshals into a space delimited string containing multiple response
// type values.
type ResponseTypeSet uint8

func (s ResponseTypeSet) UInt8() uint8 {
	return uint8(s)
}

func (s ResponseTypeSet) Contains(t ResponseType) bool {
	return s.UInt8()&t.UInt8() != 0
}

// Add adds given response types to this set and returns the updated set value. Note caller must record the returned
// value to complete the Add, or it is lost.
func (s ResponseTypeSet) Add(responseType ...ResponseType) ResponseTypeSet {
	v := s
	for _, t := range responseType {
		v = ResponseTypeSet(v.UInt8() | t.UInt8())
	}
	return v
}

// AddValues adds given response type values to this set and returns the updated set value, or an error if any value
// is not a valid response type value. Note caller must record the returned value to complete the AddValues, or it is
// lost.
func (s ResponseTypeSet) AddValues(value ...string) (ResponseTypeSet, error) {
	v := s
	for _, each := range value {
		switch each {
		case responseTypeCode:
			v = v.Add(ResponseTypeCode)
		case responseTypeToken:
			v = v.Add(ResponseTypeToken)
		case responseTypeIDToken:
			v = v.Add(ResponseTypeIDToken)
		default:
			return 0, fmt.Errorf("invalid value for spec.ResponseType: [%s]", each)
		}
	}
	return v, nil
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

func (s *ResponseTypeSet) UnmarshalJSON(bytes []byte) (err error) {
	var value string

	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return
	}

	if len(value) == 0 {
		return
	}

	*s, err = ResponseTypeSet(0).AddValues(strings.Split(value, " ")...)

	return
}
