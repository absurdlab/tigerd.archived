//go:build unit

package jose

import (
	_ "embed"
	"github.com/absurdlab/tigerd/internal/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var (
	//go:embed testdata/sign_standard_claims.jwt
	signedStandardClaimsJWT string
	//go:embed testdata/encrypt_standard_claims.jwt
	encryptedStandardClaimsJWT string
	//go:embed testdata/sign_and_encrypt_standard_claims.jwt
	signedAndEncryptedStandardClaimsJWT string
)

func TestDecode(t *testing.T) {
	jwks, err := ReadJSONWebKeySet(strings.NewReader(testJWKS))
	require.NoError(t, err)

	for _, c := range []struct {
		name  string
		token string
		opt   []DecoderOpt
	}{
		{
			name:  "decrypt signed token",
			token: signedStandardClaimsJWT,
			opt: []DecoderOpt{
				ExpectSignature(spec.RS256, jwks),
			},
		},
		{
			name:  "decrypt encrypted token",
			token: encryptedStandardClaimsJWT,
			opt: []DecoderOpt{
				ExpectEncryption(spec.RSA_OAEP_256, jwks),
			},
		},
		{
			name:  "decrypt signed and encrypted token",
			token: signedAndEncryptedStandardClaimsJWT,
			opt: []DecoderOpt{
				ExpectSignature(spec.ES256, jwks),
				ExpectEncryption(spec.RSA_OAEP_256, jwks),
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			claims := new(testStdClaims)
			err := Decode(c.token, c.opt...).Into(&claims)
			if assert.NoError(t, err) {
				claims.assert(t)
			}
		})
	}
}
