//go:build unit

package authorize_test

import (
	"github.com/absurdlab/tigerd/internal/authorize"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviderProperties_Validate(t *testing.T) {
	cases := []struct {
		name   string
		prop   *authorize.ProviderProperties
		assert func(t *testing.T, err error)
	}{
		{
			name: "correct",
			prop: &authorize.ProviderProperties{Key: "foo", Address: "localhost:30000"},
			assert: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "missing key",
			prop: &authorize.ProviderProperties{Key: "", Address: "localhost:30000"},
			assert: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "invalid address",
			prop: &authorize.ProviderProperties{Key: "foo", Address: "bar"},
			assert: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.assert(t, c.prop.Validate())
		})
	}
}
