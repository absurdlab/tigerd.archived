package jose

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"github.com/absurdlab/tigerd/internal/spec"
	"github.com/go-jose/go-jose/v3"
)

// GenerateSignatureKey generates a signature Key with the given kid and algorithm.
func GenerateSignatureKey(kid string, alg spec.SignatureAlgorithm, bits int) *JSONWebKey {
	var key any

	switch alg {
	case spec.RS256, spec.RS384, spec.RS512,
		spec.ES256, spec.ES384, spec.ES512,
		spec.PS256, spec.PS384, spec.PS512:
		key = mustGenSig(alg, bits)
	case spec.HS256:
		key = mustGenOctet(256)
	case spec.HS384:
		key = mustGenOctet(384)
	case spec.HS512:
		key = mustGenOctet(512)
	default:
		panic("unsupported algorithm")
	}

	return &JSONWebKey{
		Key:       key,
		KeyID:     kid,
		Algorithm: alg.String(),
		Use:       UseSig,
	}
}

// GenerateOctetKey generates a symmetric Key of random bytes.
func GenerateOctetKey(kid string, bytes int) *JSONWebKey {
	return &JSONWebKey{
		Key:   mustGenOctet(bytes * 8),
		KeyID: kid,
	}
}

// GenerateEncryptionKey generates an encryption Key with the given kid and algorithm.
func GenerateEncryptionKey(kid string, alg spec.EncryptionAlgorithm, bits int) *JSONWebKey {
	var key any

	switch alg {
	case spec.RSA1_5, spec.RSA_OAEP, spec.RSA_OAEP_256,
		spec.ECDH_ES, spec.ECDH_ES_A128KW, spec.ECDH_ES_A192KW, spec.ECDH_ES_A256KW:
		key = mustGenEnc(alg, bits)
	case spec.A128KW, spec.A128GCMKW:
		key = mustGenOctet(128)
	case spec.A192KW, spec.A192GCMKW:
		key = mustGenOctet(192)
	case spec.A256KW, spec.A256GCMKW:
		key = mustGenOctet(256)
	default:
		panic("unsupported algorithm")
	}

	return &JSONWebKey{
		Key:       key,
		KeyID:     kid,
		Algorithm: alg.String(),
		Use:       UseEnc,
	}
}

func mustGenSig(alg spec.SignatureAlgorithm, bits int) any {
	_, pk, err := keygenSig(jose.SignatureAlgorithm(alg.String()), bits)
	if err != nil {
		panic(err)
	}
	return pk
}

func mustGenEnc(alg spec.EncryptionAlgorithm, bits int) any {
	_, pk, err := keygenEnc(jose.KeyAlgorithm(alg.String()), bits)
	if err != nil {
		panic(err)
	}
	return pk
}

func mustGenOctet(bits int) []byte {
	if bits == 0 || bits%8 != 0 {
		panic("invalid bits: must be non-zero and divisible by 8")
	}

	key := make([]byte, bits/8)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}

	return key
}

// This method is copied directly from gopkg.in/square/go-jose.v2/jwk-keygen package.
func keygenSig(alg jose.SignatureAlgorithm, bits int) (crypto.PublicKey, crypto.PrivateKey, error) {
	switch alg {
	case jose.ES256, jose.ES384, jose.ES512, jose.EdDSA:
		keylen := map[jose.SignatureAlgorithm]int{
			jose.ES256: 256,
			jose.ES384: 384,
			jose.ES512: 521, // sic!
			jose.EdDSA: 256,
		}
		if bits != 0 && bits != keylen[alg] {
			return nil, nil, errors.New("this `alg` does not support arbitrary key length")
		}
	case jose.RS256, jose.RS384, jose.RS512, jose.PS256, jose.PS384, jose.PS512:
		if bits == 0 {
			bits = 2048
		}
		if bits < 2048 {
			return nil, nil, errors.New("too short key for RSA `alg`, 2048+ should required")
		}
	}
	switch alg {
	case jose.ES256:
		// The cryptographic operations are implemented using constant-time algorithms.
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	case jose.ES384:
		// NB: The cryptographic operations do not use constant-time algorithms.
		key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	case jose.ES512:
		// NB: The cryptographic operations do not use constant-time algorithms.
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	case jose.EdDSA:
		pub, key, err := ed25519.GenerateKey(rand.Reader)
		return pub, key, err
	case jose.RS256, jose.RS384, jose.RS512, jose.PS256, jose.PS384, jose.PS512:
		key, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	default:
		return nil, nil, errors.New("unknown `alg` for `use` = `sig`")
	}
}

// This method is copied directly from gopkg.in/square/go-jose.v2/jwk-keygen package.
func keygenEnc(alg jose.KeyAlgorithm, bits int) (crypto.PublicKey, crypto.PrivateKey, error) {
	switch alg {
	case jose.RSA1_5, jose.RSA_OAEP, jose.RSA_OAEP_256:
		if bits == 0 {
			bits = 2048
		}
		if bits < 2048 {
			return nil, nil, errors.New("too short key for RSA `alg`, 2048+ should required")
		}
		key, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	case jose.ECDH_ES, jose.ECDH_ES_A128KW, jose.ECDH_ES_A192KW, jose.ECDH_ES_A256KW:
		var crv elliptic.Curve
		switch bits {
		case 0, 256:
			crv = elliptic.P256()
		case 384:
			crv = elliptic.P384()
		case 521:
			crv = elliptic.P521()
		default:
			return nil, nil, errors.New("unknown elliptic curve bit length, use one of 256, 384, 521")
		}
		key, err := ecdsa.GenerateKey(crv, rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key.Public(), key, err
	default:
		return nil, nil, errors.New("unknown `alg` for `use` = `enc`")
	}
}
