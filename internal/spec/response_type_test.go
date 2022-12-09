package spec_test

import (
	"encoding/json"
	"github.com/absurdlab/tigerd/internal/spec"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResponseTypeSet_MarshalJSON(t *testing.T) {
	type container struct {
		Value *spec.ResponseTypeSet `json:"value,omitempty"`
	}

	cases := []struct {
		name   string
		set    *spec.ResponseTypeSet
		expect string
	}{
		{
			name:   "code",
			set:    spec.NewResponseTypeSet(spec.ResponseTypeCode),
			expect: `{"value": "code"}`,
		},
		{
			name:   "token",
			set:    spec.NewResponseTypeSet(spec.ResponseTypeToken),
			expect: `{"value": "token"}`,
		},
		{
			name:   "id_token",
			set:    spec.NewResponseTypeSet(spec.ResponseTypeIDToken),
			expect: `{"value": "id_token"}`,
		},
		{
			name:   "code token",
			set:    spec.NewResponseTypeSet(spec.ResponseTypeCode, spec.ResponseTypeToken),
			expect: `{"value": "code token"}`,
		},
		{
			name:   "code id_token",
			set:    spec.NewResponseTypeSet(spec.ResponseTypeCode, spec.ResponseTypeIDToken),
			expect: `{"value": "code id_token"}`,
		},
		{
			name:   "token id_token",
			set:    spec.NewResponseTypeSet(spec.ResponseTypeToken, spec.ResponseTypeIDToken),
			expect: `{"value": "token id_token"}`,
		},
		{
			name:   "code token id_token",
			set:    spec.NewResponseTypeSet(spec.ResponseTypeCode, spec.ResponseTypeToken, spec.ResponseTypeIDToken),
			expect: `{"value": "code token id_token"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw, err := json.Marshal(container{Value: c.set})
			assert.NoError(t, err)
			assert.JSONEq(t, c.expect, string(raw))

			var dest container
			assert.NoError(t, json.Unmarshal(raw, &dest))
			assert.Equal(t, c.set.UInt8(), dest.Value.UInt8())
		})
	}
}
