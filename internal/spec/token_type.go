package spec

import (
	"encoding/json"
	"fmt"
)

const (
	TokenTypeAccess TokenType = 1 << iota
	TokenTypeRefresh

	tokenTypeAccess  = "access_token"
	tokenTypeRefresh = "refresh_token"
)

// TokenType represents the token_type parameter in token revocation requests.
type TokenType uint8

func (t TokenType) String() string {
	switch t {
	case TokenTypeAccess:
		return tokenTypeAccess
	case TokenTypeRefresh:
		return tokenTypeRefresh
	default:
		return ""
	}
}

func (t TokenType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *TokenType) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case tokenTypeAccess:
		*t = TokenTypeAccess
	case tokenTypeRefresh:
		*t = TokenTypeRefresh
	default:
		return fmt.Errorf("invalid value for spec.TokenType [%s]", value)
	}

	return nil
}
