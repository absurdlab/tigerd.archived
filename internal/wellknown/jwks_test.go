package wellknown_test

import (
	"github.com/absurdlab/tigerd/internal/wellknown"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJSONWebKeySet(t *testing.T) {
	jwks, err := wellknown.NewJSONWebKeySet(&wellknown.JSONWebKeySetProperties{File: "testdata/jwks.json"})
	if assert.NoError(t, err) {
		assert.True(t, jwks.Size() > 0)
	}
}
