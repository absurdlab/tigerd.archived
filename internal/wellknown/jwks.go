package wellknown

import (
	"errors"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/absurdlab/tigerd/internal/jose"
	"github.com/absurdlab/tigerd/internal/spec"
	"io"
	"os"
	"strings"
)

var (
	ErrJWKS = errors.New("invalid jwks")
)

// JSONWebKeySetProperties is the configuration properties for reading server's JSON Web Key Set. The key set can be
// read from either a File or an Inline string. When both option are specified, the File option takes precedence.
type JSONWebKeySetProperties struct {
	File   string
	Inline string
}

// NewJSONWebKeySet reads a jose.JSONWebKeySet by means specified in JSONWebKeySetProperties.
func NewJSONWebKeySet(props *JSONWebKeySetProperties) (*jose.JSONWebKeySet, error) {
	var (
		reader io.Reader
		err    error
	)
	switch {
	case len(props.File) > 0:
		reader, err = os.Open(props.File)
		if err != nil {
			return nil, fault.Wrap(ErrJWKS,
				ftag.With(spec.ErrKindInvalidRequest),
				fmsg.WithDesc(err.Error(), "Failed to open jwks file."),
			)
		}
	case len(props.Inline) > 0:
		reader, err = strings.NewReader(props.Inline), nil
	default:
		err = fault.Wrap(ErrJWKS,
			ftag.With(spec.ErrKindInvalidRequest),
			fmsg.WithDesc(err.Error(), "Expect either file or inline option for jwks."),
		)
	}
	if err != nil {
		return nil, err
	}

	jwks, err := jose.ReadJSONWebKeySet(reader)
	if err != nil {
		return nil, fault.Wrap(ErrJWKS,
			ftag.With(spec.ErrKindInvalidRequest),
			fmsg.WithDesc(err.Error(), "Invalid jwks definition."),
		)
	}

	return jwks, nil
}
