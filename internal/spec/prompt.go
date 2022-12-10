package spec

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	PromptLogin Prompt = 1 << iota
	PromptSelectAccount
	PromptConsent
	PromptNone

	promptLogin         = "login"
	promptSelectAccount = "select_account"
	promptConsent       = "consent"
	promptNone          = "none"
)

// Prompt presents the prompt parameter defined in OpenID Connect 1.0.
type Prompt uint8

func (p Prompt) UInt8() uint8 {
	return uint8(p)
}

func (p Prompt) String() string {
	switch p {
	case PromptLogin:
		return promptLogin
	case PromptSelectAccount:
		return promptSelectAccount
	case PromptConsent:
		return promptConsent
	case PromptNone:
		return promptNone
	default:
		return ""
	}
}

// PromptSet is a bitmask of zero or more Prompt that models multi prompt format. It is distinct from an array of
// Prompt. Instead, it marshals into a space delimited string containing multiple prompt values.
type PromptSet uint8

func (s PromptSet) UInt8() uint8 {
	return uint8(s)
}

func (s PromptSet) Contains(t Prompt) bool {
	return s.UInt8()&t.UInt8() != 0
}

// IsValid returns true if this PromptSet contains a valid combination of Prompt values.
func (s PromptSet) IsValid() bool {
	switch {
	// "none" cannot co-exist with other prompt - since PromptNone is enumerated last,
	// there cannot be a greater value in this bitmask.
	case s.UInt8() > PromptNone.UInt8():
		return false

	case s.Contains(PromptLogin) && s.Contains(PromptSelectAccount):
		return false

	default:
		return true
	}
}

// Add adds given prompt to this set and returns the updated set value. Note caller must record the returned
// value to complete the Add, or it is lost.
func (s PromptSet) Add(prompt ...Prompt) PromptSet {
	v := s
	for _, t := range prompt {
		v = PromptSet(v.UInt8() | t.UInt8())
	}
	return v
}

// AddValues adds given prompt to this set and returns the updated set value, or an error if any value
// is not a valid prompt. Note caller must record the returned value to complete the AddValues, or it is
// lost.
func (s PromptSet) AddValues(value ...string) (PromptSet, error) {
	v := s
	for _, each := range value {
		switch each {
		case promptLogin:
			v = v.Add(PromptLogin)
		case promptSelectAccount:
			v = v.Add(PromptSelectAccount)
		case promptConsent:
			v = v.Add(PromptConsent)
		case promptNone:
			v = v.Add(PromptNone)
		default:
			return 0, fmt.Errorf("invalid value for spec.Prompt: [%s]", each)
		}
	}
	return v, nil
}

func (s PromptSet) MarshalJSON() ([]byte, error) {
	if s == 0 {
		return nil, nil
	}

	var values []string
	for _, t := range []Prompt{PromptLogin, PromptSelectAccount, PromptConsent, PromptNone} {
		if s.Contains(t) {
			values = append(values, t.String())
		}
	}

	return json.Marshal(strings.Join(values, " "))
}

func (s *PromptSet) UnmarshalJSON(bytes []byte) (err error) {
	var value string

	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return
	}

	if len(value) == 0 {
		return
	}

	*s, err = PromptSet(0).AddValues(strings.Split(value, " ")...)

	return
}
