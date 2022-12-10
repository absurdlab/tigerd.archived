package spec

import (
	"encoding/json"
	"fmt"
)

const (
	GrantTypeAuthorizationCode GrantType = 1 << iota
	GrantTypeImplicit
	GrantTypePassword
	GrantTypeClientCredentials
	GrantTypeRefreshToken

	grantTypeAuthorizationCode = "authorization_code"
	grantTypeImplicit          = "implicit"
	grantTypePassword          = "password"
	grantTypeClientCredentials = "client_credentials"
	grantTypeRefreshToken      = "refresh_token"
)

type GrantType uint8

func (g GrantType) String() string {
	switch g {
	case GrantTypeAuthorizationCode:
		return grantTypeAuthorizationCode
	case GrantTypeImplicit:
		return grantTypeImplicit
	case GrantTypePassword:
		return grantTypePassword
	case GrantTypeClientCredentials:
		return grantTypeClientCredentials
	case GrantTypeRefreshToken:
		return grantTypeRefreshToken
	default:
		return ""
	}
}

func (g GrantType) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

func (g *GrantType) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case grantTypeAuthorizationCode:
		*g = GrantTypeAuthorizationCode
	case grantTypeImplicit:
		*g = GrantTypeImplicit
	case grantTypePassword:
		*g = GrantTypePassword
	case grantTypeClientCredentials:
		*g = GrantTypeClientCredentials
	case grantTypeRefreshToken:
		*g = GrantTypeRefreshToken
	default:
		return fmt.Errorf("invalid spec.GrantType value [%s]", value)
	}

	return nil
}
