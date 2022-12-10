package spec

import (
	"encoding/json"
	"fmt"
	"github.com/go-jose/go-jose/v3"
)

const (
	NoSignature SignatureAlgorithm = 1 << iota
	HS256
	HS384
	HS512
	RS256
	RS384
	RS512
	ES256
	ES384
	ES512
	PS256
	PS384
	PS512
)

// SignatureAlgorithm is the cryptographic algorithms defined in JSON Web Algorithms used to perform signing and verification.
type SignatureAlgorithm uint16

func (g SignatureAlgorithm) IsNoneOrEmpty() bool {
	return g == 0 || g == NoSignature
}

func (g SignatureAlgorithm) String() string {
	switch g {
	case NoSignature:
		return "none"
	case HS256:
		return string(jose.HS256)
	case HS384:
		return string(jose.HS384)
	case HS512:
		return string(jose.HS512)
	case RS256:
		return string(jose.RS256)
	case RS384:
		return string(jose.RS384)
	case RS512:
		return string(jose.RS512)
	case ES256:
		return string(jose.ES256)
	case ES384:
		return string(jose.ES384)
	case ES512:
		return string(jose.ES512)
	case PS256:
		return string(jose.PS256)
	case PS384:
		return string(jose.PS384)
	case PS512:
		return string(jose.PS512)
	default:
		return ""
	}
}

func (g SignatureAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

func (g *SignatureAlgorithm) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case "none":
		*g = NoSignature
	case string(jose.HS256):
		*g = HS256
	case string(jose.HS384):
		*g = HS384
	case string(jose.HS512):
		*g = HS512
	case string(jose.RS256):
		*g = RS256
	case string(jose.RS384):
		*g = RS384
	case string(jose.RS512):
		*g = RS512
	case string(jose.ES256):
		*g = ES256
	case string(jose.ES384):
		*g = ES384
	case string(jose.ES512):
		*g = ES512
	case string(jose.PS256):
		*g = PS256
	case string(jose.PS384):
		*g = PS384
	case string(jose.PS512):
		*g = PS512
	default:
		return fmt.Errorf("invalid spec.SignatureAlgorithm value [%s]", value)
	}

	return nil
}

//goland:noinspection GoSnakeCaseUsage
const (
	NoEncryption EncryptionAlgorithm = 1 << iota
	ED25519
	RSA1_5
	RSA_OAEP
	RSA_OAEP_256
	A128KW
	A192KW
	A256KW
	DIRECT
	ECDH_ES
	ECDH_ES_A128KW
	ECDH_ES_A192KW
	ECDH_ES_A256KW
	A128GCMKW
	A192GCMKW
	A256GCMKW
	PBES2_HS256_A128KW
	PBES2_HS384_A192KW
	PBES2_HS512_A256KW
)

// EncryptionAlgorithm is the cryptographic algorithm defined in JSON Web Algorithm to perform encryption and decryption.
type EncryptionAlgorithm uint32

func (g EncryptionAlgorithm) IsNoneOrEmpty() bool {
	return g == 0 || g == NoEncryption
}

func (g EncryptionAlgorithm) String() string {
	switch g {
	case NoEncryption:
		return "none"
	case ED25519:
		return string(jose.ED25519)
	case RSA1_5:
		return string(jose.RSA1_5)
	case RSA_OAEP:
		return string(jose.RSA_OAEP)
	case RSA_OAEP_256:
		return string(jose.RSA_OAEP_256)
	case A128KW:
		return string(jose.A128KW)
	case A192KW:
		return string(jose.A192KW)
	case A256KW:
		return string(jose.A256KW)
	case DIRECT:
		return string(jose.DIRECT)
	case ECDH_ES:
		return string(jose.ECDH_ES)
	case ECDH_ES_A128KW:
		return string(jose.ECDH_ES_A128KW)
	case ECDH_ES_A192KW:
		return string(jose.ECDH_ES_A192KW)
	case ECDH_ES_A256KW:
		return string(jose.ECDH_ES_A256KW)
	case A128GCMKW:
		return string(jose.A128GCMKW)
	case A192GCMKW:
		return string(jose.A192GCMKW)
	case A256GCMKW:
		return string(jose.A256GCMKW)
	case PBES2_HS256_A128KW:
		return string(jose.PBES2_HS256_A128KW)
	case PBES2_HS384_A192KW:
		return string(jose.PBES2_HS384_A192KW)
	case PBES2_HS512_A256KW:
		return string(jose.PBES2_HS512_A256KW)
	default:
		return ""
	}
}

func (g EncryptionAlgorithm) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

func (g *EncryptionAlgorithm) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case "none":
		*g = NoEncryption
	case string(jose.ED25519):
		*g = ED25519
	case string(jose.RSA1_5):
		*g = RSA1_5
	case string(jose.RSA_OAEP):
		*g = RSA_OAEP
	case string(jose.RSA_OAEP_256):
		*g = RSA_OAEP_256
	case string(jose.A128KW):
		*g = A128KW
	case string(jose.A192KW):
		*g = A192KW
	case string(jose.A256KW):
		*g = A256KW
	case string(jose.DIRECT):
		*g = DIRECT
	case string(jose.ECDH_ES):
		*g = ECDH_ES
	case string(jose.ECDH_ES_A128KW):
		*g = ECDH_ES_A128KW
	case string(jose.ECDH_ES_A192KW):
		*g = ECDH_ES_A192KW
	case string(jose.ECDH_ES_A256KW):
		*g = ECDH_ES_A256KW
	case string(jose.A128GCMKW):
		*g = A128GCMKW
	case string(jose.A192GCMKW):
		*g = A192GCMKW
	case string(jose.A256GCMKW):
		*g = A256GCMKW
	case string(jose.PBES2_HS256_A128KW):
		*g = PBES2_HS256_A128KW
	case string(jose.PBES2_HS384_A192KW):
		*g = PBES2_HS384_A192KW
	case string(jose.PBES2_HS512_A256KW):
		*g = PBES2_HS512_A256KW
	default:
		return fmt.Errorf("invalid spec.EncryptionAlgorithm value [%s]", value)
	}

	return nil
}

//goland:noinspection GoSnakeCaseUsage
const (
	NoEncoding EncryptionEncoding = 1 << iota
	A128CBC_HS256
	A192CBC_HS384
	A256CBC_HS512
	A128GCM
	A192GCM
	A256GCM
)

// EncryptionEncoding is the cryptographic algorithm defined in JSON Web Algorithm to perform encoding and decoding.
type EncryptionEncoding uint8

func (g EncryptionEncoding) IsNoneOrEmpty() bool {
	return g == 0 || g == NoEncoding
}

func (g EncryptionEncoding) String() string {
	switch g {
	case NoEncoding:
		return "none"
	case A128CBC_HS256:
		return string(jose.A128CBC_HS256)
	case A192CBC_HS384:
		return string(jose.A192CBC_HS384)
	case A256CBC_HS512:
		return string(jose.A256CBC_HS512)
	case A128GCM:
		return string(jose.A128GCM)
	case A192GCM:
		return string(jose.A192GCM)
	case A256GCM:
		return string(jose.A256GCM)
	default:
		return ""
	}
}

func (g EncryptionEncoding) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.String())
}

func (g *EncryptionEncoding) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case "none":
		*g = NoEncoding
	case string(jose.A128CBC_HS256):
		*g = A128CBC_HS256
	case string(jose.A192CBC_HS384):
		*g = A192CBC_HS384
	case string(jose.A256CBC_HS512):
		*g = A256CBC_HS512
	case string(jose.A128GCM):
		*g = A128GCM
	case string(jose.A192GCM):
		*g = A192GCM
	case string(jose.A256GCM):
		*g = A256GCM
	default:
		return fmt.Errorf("invalid spec.EncryptionEncoding value [%s]", value)
	}

	return nil
}
