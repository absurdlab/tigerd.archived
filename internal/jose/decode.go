package jose

import (
	"errors"
	"fmt"
	"github.com/absurdlab/tigerd/internal/spec"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

var (
	// ErrDecode should the base marker error for jose decoding.
	ErrDecode = errors.New("jose decode error")
)

// Decode performs JWT/JWE decode operations on the supplied token value. Use ExpectSignature and/or ExpectEncryption to
// specify the expected shape of the token. This setup function returns a Decoder which has a single public method
// Decoder.Into to decode the token into one or more claims structure.
func Decode(token string, opts ...DecoderOpt) *Decoder {
	n := new(Decoder)
	n.token = token
	for _, opt := range opts {
		opt(n)
	}
	return n
}

// DecoderOpt should the option for Decode.
type DecoderOpt func(n *Decoder)

// PeekOnly instructs the Decoder to skip signature verification and decode the payload directly.
func PeekOnly() DecoderOpt {
	return func(n *Decoder) {
		n.peekMode = true
	}
}

// ExpectSignature instructs the Decoder to expect the token to have been signed and therefore perform signature verification.
func ExpectSignature(alg spec.SignatureAlgorithm, jwks *JSONWebKeySet) DecoderOpt {
	return func(n *Decoder) {
		if alg.IsNoneOrEmpty() || jwks == nil {
			return
		}

		n.verificationAlg = alg
		n.verificationKeySet = jwks
	}
}

// ExpectEncryption instructs the Decoder to expect the token to have been encrypted and therefore perform decryption.
func ExpectEncryption(alg spec.EncryptionAlgorithm, jwks *JSONWebKeySet) DecoderOpt {
	return func(n *Decoder) {
		if alg.IsNoneOrEmpty() || jwks == nil {
			return
		}

		n.decryptionAlg = alg
		n.decryptionKeySet = jwks
	}
}

type Decoder struct {
	token              string
	verificationAlg    spec.SignatureAlgorithm
	verificationKeySet *JSONWebKeySet
	decryptionAlg      spec.EncryptionAlgorithm
	decryptionKeySet   *JSONWebKeySet
	peekMode           bool
}

// Into decodes the token into the supplied destination structures.
func (n *Decoder) Into(dest ...any) (err error) {
	defer func() {
		if err != nil && !errors.Is(err, ErrDecode) {
			err = fmt.Errorf("%w: %s", ErrDecode, err)
		}
	}()

	if len(n.token) == 0 {
		return errors.New("empty token")
	}

	if len(dest) == 0 {
		return errors.New("no destination")
	}

	if n.peekMode {
		return n.peek(dest...)
	}

	doVerify, doDecrypt := n.expectSignature(), n.expectEncryption()

	switch {
	case !doVerify && !doDecrypt:
		return errors.New("invalid configuration")
	case doVerify && doDecrypt:
		return n.decryptAndVerify(dest...)
	case doVerify:
		return n.verify(dest...)
	case doDecrypt:
		return n.decrypt(dest...)
	default:
		panic("impossible case")
	}
}

func (n *Decoder) peek(dest ...any) error {
	jsonWebToken, err := jwt.ParseSigned(n.token)
	if err != nil {
		return err
	}

	return jsonWebToken.UnsafeClaimsWithoutVerification(dest...)
}

func (n *Decoder) verify(dest ...any) error {
	jsonWebToken, err := jwt.ParseSigned(n.token)
	if err != nil {
		return err
	}

	jsonWebKey, err := n.extractVerificationKey(jsonWebToken.Headers)
	if err != nil {
		return err
	}

	return jsonWebToken.Claims(jsonWebKey.Public().Key, dest...)
}

func (n *Decoder) decrypt(dest ...any) error {
	jsonWebToken, err := jwt.ParseEncrypted(n.token)
	if err != nil {
		return err
	}

	jsonWebKey, err := n.extractDecryptionKey(jsonWebToken.Headers)
	if err != nil {
		return err
	}

	return jsonWebToken.Claims(jsonWebKey.Key, dest...)
}

func (n *Decoder) decryptAndVerify(dest ...any) error {
	nestedJWT, err := jwt.ParseSignedAndEncrypted(n.token)
	if err != nil {
		return err
	}

	decryptJsonWebKey, err := n.extractDecryptionKey(nestedJWT.Headers)
	if err != nil {
		return err
	}

	jsonWebToken, err := nestedJWT.Decrypt(decryptJsonWebKey.Key)
	if err != nil {
		return err
	}

	verifyJsonWebKey, err := n.extractVerificationKey(jsonWebToken.Headers)
	if err != nil {
		return err
	}

	return jsonWebToken.Claims(verifyJsonWebKey.Public().Key, dest...)
}

func (n *Decoder) expectSignature() bool {
	return !n.verificationAlg.IsNoneOrEmpty() && n.verificationKeySet != nil
}

func (n *Decoder) expectEncryption() bool {
	return !n.decryptionAlg.IsNoneOrEmpty() && n.decryptionKeySet != nil
}

func (n *Decoder) extractVerificationKey(headers []jose.Header) (*JSONWebKey, error) {
	kid, alg := getKeyIdFromHeaders(headers), getAlgFromHeaders(headers)
	if len(kid) == 0 {
		return nil, errors.New("jws missing kid header")
	} else if len(alg) == 0 {
		return nil, errors.New("jws missing alg header")
	}

	jsonWebKey := n.verificationKeySet.FindByKeyID(kid)
	if jsonWebKey == nil {
		return nil, fmt.Errorf("no key by kid %s", kid)
	}

	switch {
	case alg != jsonWebKey.Algorithm:
		return nil, errors.New("jws token alg header mismatch with key algorithm")
	case alg != n.verificationAlg.String():
		return nil, errors.New("verify algorithm mismatch")
	}

	return jsonWebKey, nil
}

func (n *Decoder) extractDecryptionKey(headers []jose.Header) (*JSONWebKey, error) {
	kid, alg := getKeyIdFromHeaders(headers), getAlgFromHeaders(headers)
	if len(kid) == 0 {
		return nil, errors.New("jwe missing kid header")
	} else if len(alg) == 0 {
		return nil, errors.New("jwe missing alg header")
	}

	jsonWebKey := n.decryptionKeySet.FindByKeyID(kid)
	if jsonWebKey == nil {
		return nil, fmt.Errorf("no key by kid %s", kid)
	}

	switch {
	case alg != jsonWebKey.Algorithm:
		return nil, errors.New("jwe token alg header mismatch with key algorithm")
	case alg != n.decryptionAlg.String():
		return nil, errors.New("decryption algorithm mismatch")
	}

	return jsonWebKey, nil
}

func getKeyIdFromHeaders(headers []jose.Header) string {
	for _, header := range headers {
		if kid := header.KeyID; len(kid) > 0 {
			return kid
		}
	}
	return ""
}

func getAlgFromHeaders(headers []jose.Header) string {
	for _, header := range headers {
		if alg := header.Algorithm; len(alg) > 0 {
			return alg
		}
	}
	return ""
}
