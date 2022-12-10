package spec_test

import (
	"encoding/json"
	"fmt"
	"github.com/absurdlab/tigerd/internal/spec"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPromptSet_MarshalJSON(t *testing.T) {
	type container struct {
		Value spec.PromptSet `json:"value,omitempty"`
	}

	cases := []struct {
		name   string
		set    spec.PromptSet
		expect string
	}{
		{
			name:   "login",
			set:    spec.PromptSet(0).Add(spec.PromptLogin),
			expect: `{"value": "login"}`,
		},
		{
			name:   "select_account",
			set:    spec.PromptSet(0).Add(spec.PromptSelectAccount),
			expect: `{"value": "select_account"}`,
		},
		{
			name:   "consent",
			set:    spec.PromptSet(0).Add(spec.PromptConsent),
			expect: `{"value": "consent"}`,
		},
		{
			name:   "none",
			set:    spec.PromptSet(0).Add(spec.PromptNone),
			expect: `{"value": "none"}`,
		},
		{
			name:   "login consent",
			set:    spec.PromptSet(0).Add(spec.PromptLogin, spec.PromptConsent),
			expect: `{"value": "login consent"}`,
		},
		{
			name:   "select_account consent",
			set:    spec.PromptSet(0).Add(spec.PromptSelectAccount, spec.PromptConsent),
			expect: `{"value": "select_account consent"}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw, err := json.Marshal(container{Value: c.set})
			assert.NoError(t, err)
			assert.JSONEq(t, c.expect, string(raw))

			var dest container
			assert.NoError(t, json.Unmarshal(raw, &dest))
			assert.Equal(t, c.set, dest.Value)
		})
	}
}

func TestPromptSet_IsValid(t *testing.T) {
	for i, c := range []struct {
		set   spec.PromptSet
		valid bool
	}{
		{set: spec.PromptSet(0).Add(spec.PromptLogin), valid: true},
		{set: spec.PromptSet(0).Add(spec.PromptSelectAccount), valid: true},
		{set: spec.PromptSet(0).Add(spec.PromptConsent), valid: true},
		{set: spec.PromptSet(0).Add(spec.PromptNone), valid: true},
		{set: spec.PromptSet(0).Add(spec.PromptLogin, spec.PromptConsent), valid: true},
		{set: spec.PromptSet(0).Add(spec.PromptSelectAccount, spec.PromptConsent), valid: true},
		{set: spec.PromptSet(0).Add(spec.PromptLogin, spec.PromptSelectAccount), valid: false},
		{set: spec.PromptSet(0).Add(spec.PromptLogin, spec.PromptNone), valid: false},
		{set: spec.PromptSet(0).Add(spec.PromptSelectAccount, spec.PromptNone), valid: false},
		{set: spec.PromptSet(0).Add(spec.PromptConsent, spec.PromptNone), valid: false},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, c.valid, c.set.IsValid())
		})
	}
}
