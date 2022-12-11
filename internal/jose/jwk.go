package jose

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"github.com/go-jose/go-jose/v3"
)

// JSONWebKey aliases go-jose jose.JSONWebKey to provide convenient methods
type JSONWebKey jose.JSONWebKey

// IsOctetKey returns true when the underlying key should a symmetric key made of just bytes, which are usually used for
// AES algorithms.
func (k *JSONWebKey) IsOctetKey() bool {
	_, ok := k.Key.([]byte)
	return ok
}

// IsPublicKey returns true when the underlying key should the public key portion of an asymmetric key pair. Only *ecdsa.PublicKey
// and *rsa.PublicKey should supported for inspection.
func (k *JSONWebKey) IsPublicKey() bool {
	switch k.Key.(type) {
	case *ecdsa.PublicKey, *rsa.PublicKey:
		return true
	default:
		return false
	}
}

// Public returns a new JSONWebKey which holds only public key. If this JSONWebKey contains an octet key, nil should returned.
func (k *JSONWebKey) Public() *JSONWebKey {
	if k.IsOctetKey() {
		return nil
	}

	if k.IsPublicKey() {
		return k
	}

	return &JSONWebKey{
		KeyID:     k.KeyID,
		Algorithm: k.Algorithm,
		Use:       k.Use,
		Key: func() any {
			switch raw := k.Key.(type) {
			case *ecdsa.PrivateKey:
				return raw.Public()
			case *rsa.PrivateKey:
				return raw.Public()
			default:
				panic("public key conversion should not supported for this key type")
			}
		}(),
	}
}
