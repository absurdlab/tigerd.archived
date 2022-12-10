package spec

import (
	"encoding/json"
	"fmt"
)

const (
	NoAuthenticationMethod AuthenticationMethod = 1 << iota
	ClientSecretBasic
	ClientSecretPost
	ClientSecretJWT
	PrivateKeyJWT

	authMethodNone    = "none"
	clientSecretBasic = "client_secret_basic"
	clientSecretPost  = "client_secret_post"
	clientSecretJWT   = "client_secret_jwt"
	privateKeyJWT     = "private_key_jwt"
)

// AuthenticationMethod represents the values of token_endpoint_auth_method specified in OAuth 2.0 and OpenID Connect 1.0.
type AuthenticationMethod uint8

func (m AuthenticationMethod) String() string {
	switch m {
	case NoAuthenticationMethod:
		return authMethodNone
	case ClientSecretBasic:
		return clientSecretBasic
	case ClientSecretPost:
		return clientSecretPost
	case ClientSecretJWT:
		return clientSecretJWT
	case PrivateKeyJWT:
		return privateKeyJWT
	default:
		return ""
	}
}

func (m AuthenticationMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *AuthenticationMethod) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case authMethodNone:
		*m = NoAuthenticationMethod
	case clientSecretBasic:
		*m = ClientSecretBasic
	case clientSecretPost:
		*m = ClientSecretPost
	case clientSecretJWT:
		*m = ClientSecretJWT
	case privateKeyJWT:
		*m = PrivateKeyJWT
	default:
		return fmt.Errorf("invalid value for spec.AuthenticationMethod [%s]", value)
	}

	return nil
}
