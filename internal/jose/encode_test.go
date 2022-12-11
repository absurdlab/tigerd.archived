package jose

import (
	_ "embed"
	"github.com/absurdlab/tigerd/internal/spec"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

const (
	fixedTimestamp int64 = 1664586927
)

var (
	//go:embed testdata/jwks.json
	testJWKS string
)

func TestEncode(t *testing.T) {
	jwks, err := ReadJSONWebKeySet(strings.NewReader(testJWKS))
	require.NoError(t, err)

	for _, c := range []struct {
		name   string
		claims any
		opt    []EncoderOpt
		assert func(t *testing.T, token string)
	}{
		{
			name:   "sign standard claims",
			claims: new(testStdClaims).init(),
			opt: []EncoderOpt{
				WithSignature(spec.RS256, jwks),
			},
			assert: func(t *testing.T, token string) {
				parsed, err := jwt.ParseSigned(token)
				if assert.NoError(t, err) {
					claims := new(testStdClaims)
					_ = parsed.UnsafeClaimsWithoutVerification(&claims)
					claims.assert(t)
				}
			},
		},
		{
			name:   "encrypt standard claims",
			claims: new(testStdClaims).init(),
			opt: []EncoderOpt{
				WithEncryption(spec.RSA_OAEP_256, spec.A128GCM, jwks),
			},
			assert: func(t *testing.T, token string) {
				parsed, err := jwt.ParseEncrypted(token)
				if assert.NoError(t, err) {
					claims := new(testStdClaims)
					_ = parsed.Claims(jwks.FindForEncryption(spec.RSA_OAEP_256).Key, &claims)
					claims.assert(t)
				}
			},
		},
		{
			name:   "sign and encrypt standard claims",
			claims: new(testStdClaims).init(),
			opt: []EncoderOpt{
				WithSignature(spec.ES256, jwks),
				WithEncryption(spec.RSA_OAEP_256, spec.A128GCM, jwks),
			},
			assert: func(t *testing.T, token string) {
				parsed, err := jwt.ParseSignedAndEncrypted(token)
				if assert.NoError(t, err) {
					parsed2, err := parsed.Decrypt(jwks.FindForEncryption(spec.RSA_OAEP_256).Key)
					if assert.NoError(t, err) {
						claims := new(testStdClaims)
						_ = parsed2.UnsafeClaimsWithoutVerification(&claims)
						claims.assert(t)
					}
				}
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			encoded, err := Encode(c.claims, c.opt...)
			if assert.NoError(t, err) {
				t.Log(encoded)
				c.assert(t, encoded)
			}
		})
	}
}

type testStdClaims struct {
	*StdClaims
}

func (c *testStdClaims) init() *testStdClaims {
	c.StdClaims = new(StdClaims).
		WithID("jti").
		WithIssuer("iss").
		WithAudience("aud").
		WithSubject("sub").
		WithIssuedAt(time.Unix(fixedTimestamp, 0)).
		WithExpiry(time.Unix(fixedTimestamp, 0)).
		WithNotBefore(time.Unix(fixedTimestamp, 0))
	return c
}

func (c *testStdClaims) assert(t *testing.T) {
	assert.Equal(t, "jti", c.ID)
	assert.Equal(t, "iss", c.Issuer)
	assert.Equal(t, "aud", c.Audience[0])
	assert.Equal(t, "sub", c.Subject)
	assert.Equal(t, fixedTimestamp, c.IssuedAt.Time().Unix())
	assert.Equal(t, fixedTimestamp, c.Expiry.Time().Unix())
	assert.Equal(t, fixedTimestamp, c.NotBefore.Time().Unix())
}
