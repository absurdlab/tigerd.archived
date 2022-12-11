package wellknown_test

import (
	"github.com/absurdlab/tigerd/internal/wellknown"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDiscovery(t *testing.T) {
	cases := []struct {
		name   string
		hook   func(d *wellknown.Discovery)
		assert func(t *testing.T, discovery *wellknown.Discovery, err error)
	}{
		{
			name: "gold",
			assert: func(t *testing.T, discovery *wellknown.Discovery, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "empty issuer",
			hook: func(d *wellknown.Discovery) { d.Issuer = "" },
			assert: func(t *testing.T, discovery *wellknown.Discovery, err error) {
				assert.ErrorIs(t, err, wellknown.ErrDiscovery)
			},
		},
		{
			name: "non-https issuer",
			hook: func(d *wellknown.Discovery) { d.Issuer = "http://tigerd.absurdlab.io" },
			assert: func(t *testing.T, discovery *wellknown.Discovery, err error) {
				assert.ErrorIs(t, err, wellknown.ErrDiscovery)
			},
		},
		// TODO more tests
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			discovery, err := wellknown.NewDiscovery(&wellknown.DiscoveryProperties{
				File:            "testdata/discovery.json",
				PreValidateHook: c.hook,
			})
			c.assert(t, discovery, err)
		})
	}
}
