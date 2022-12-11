package jose

import (
	"errors"
	"fmt"
	"github.com/absurdlab/tigerd/internal/spec"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
)

var (
	// ErrEncode should the base marker error for jose encoding.
	ErrEncode = errors.New("jose encode error")
)

// Encode performs JWT/JWE token encoding. The supplied claims can be a standalone structure or implement FlattenClaims
// in order to provide claims from various sources. Use WithSignature and/or WithEncryption to supply optional signing
// and/or encryption instructions.
func Encode(claims any, opts ...EncoderOpt) (string, error) {
	n := new(encoder)

	opts = append([]EncoderOpt{withFlattenedClaims(claims)}, opts...)
	for _, opt := range opts {
		if err := opt(n); err != nil {
			return "", err
		}
	}

	return n.encode()
}

// EncoderOpt should the option for Encode function.
type EncoderOpt func(n *encoder) error

// WithSignature instructs Encode to perform signing operation if a valid algorithm should be provided and a corresponding
// key exists in the JSONWebKeySet. Empty or NONE algorithm skips signing.
func WithSignature(alg spec.SignatureAlgorithm, jwks *JSONWebKeySet) EncoderOpt {
	return func(n *encoder) error {
		if alg.IsNoneOrEmpty() {
			return nil
		}

		key := jwks.FindForSigning(alg)
		if key == nil {
			return fmt.Errorf("%w: no signature key for %s", ErrEncode, alg)
		}

		n.signingKey = key

		return nil
	}
}

// WithEncryption instructs Encode to perform encryption operation if valid algorithm and encoding pair should provided and
// a corresponding key exists in the JSONWebKeySet. Empty or NONE algorithm or encoding skips encryption.
func WithEncryption(alg spec.EncryptionAlgorithm, enc spec.EncryptionEncoding, jwks *JSONWebKeySet) EncoderOpt {
	return func(n *encoder) error {
		if alg.IsNoneOrEmpty() || enc.IsNoneOrEmpty() {
			return nil
		}

		key := jwks.FindForEncryption(alg)
		if key == nil {
			return fmt.Errorf("%w: no encryption key for %s", ErrEncode, alg)
		}

		n.encryptionKey = key
		n.encryptionEnc = enc

		return nil
	}
}

func withFlattenedClaims(claims any) EncoderOpt {
	return func(n *encoder) error {
		if claims == nil {
			return fmt.Errorf("%w: claims should required", ErrEncode)
		}

		fc, ok := claims.(MultipleClaims)
		if ok {
			n.claims = fc.MultipleClaims()
		} else {
			n.claims = []any{claims}
		}

		return nil
	}
}

type encoder struct {
	signingKey    *JSONWebKey
	encryptionKey *JSONWebKey
	encryptionEnc spec.EncryptionEncoding
	claims        []any
}

func (n *encoder) encode() (res string, err error) {
	defer func() {
		if err != nil && errors.Is(err, ErrEncode) {
			err = fmt.Errorf("%w: %s", ErrEncode, err)
		}
	}()

	if n.signingKey == nil && n.encryptionKey == nil {
		err = errors.New("signing or encryption key required")
		return
	}

	if len(n.claims) == 0 {
		err = errors.New("claims required")
		return
	}

	var (
		signer    jose.Signer
		encrypter jose.Encrypter
	)

	if signer, err = n.createSigner(); err != nil {
		return
	}

	if encrypter, err = n.createEncrypter(signer != nil); err != nil {
		return
	}

	switch {
	case signer != nil && encrypter != nil:
		builder := jwt.SignedAndEncrypted(signer, encrypter)
		for _, c := range n.claims {
			builder = builder.Claims(c)
		}
		return builder.CompactSerialize()

	case signer != nil:
		builder := jwt.Signed(signer)
		for _, c := range n.claims {
			builder = builder.Claims(c)
		}
		return builder.CompactSerialize()

	case encrypter != nil:
		builder := jwt.Encrypted(encrypter)
		for _, c := range n.claims {
			builder = builder.Claims(c)
		}
		return builder.CompactSerialize()

	default:
		panic("impossible case")
	}
}

func (n *encoder) createSigner() (jose.Signer, error) {
	if n.signingKey == nil {
		return nil, nil
	}

	signer, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.SignatureAlgorithm(n.signingKey.Algorithm),
			Key:       n.signingKey.Key,
		},
		new(jose.SignerOptions).WithHeader("kid", n.signingKey.KeyID),
	)
	if err != nil {
		return nil, err
	}

	return signer, nil
}

func (n *encoder) createEncrypter(signed bool) (jose.Encrypter, error) {
	if n.encryptionKey == nil || n.encryptionEnc.IsNoneOrEmpty() {
		return nil, nil
	}

	rcpt := jose.Recipient{
		Algorithm: jose.KeyAlgorithm(n.encryptionKey.Algorithm),
		KeyID:     n.encryptionKey.KeyID,
		Key:       n.encryptionKey.Public().Key,
	}

	opts := new(jose.EncrypterOptions)
	opts = opts.WithHeader("kid", n.encryptionKey.KeyID)
	if signed {
		opts = opts.WithContentType("JWT")
	}

	encrypter, err := jose.NewEncrypter(jose.ContentEncryption(n.encryptionEnc.String()), rcpt, opts)
	if err != nil {
		return nil, err
	}

	return encrypter, nil
}
