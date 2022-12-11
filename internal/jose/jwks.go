package jose

import (
	"encoding/json"
	"github.com/absurdlab/tigerd/internal/spec"
	"github.com/go-jose/go-jose/v3"
	"io"
	"time"
)

// NewJSONWebKeySet creates a new key set with the given keys.
func NewJSONWebKeySet(keys ...*JSONWebKey) *JSONWebKeySet {
	s := &JSONWebKeySet{keys: map[string]*JSONWebKey{}}
	for _, key := range keys {
		s.keys[key.KeyID] = key
	}
	return s
}

// ReadJSONWebKeySet create new KeySet with data from the reader
func ReadJSONWebKeySet(reader io.Reader) (*JSONWebKeySet, error) {
	jwks := new(jose.JSONWebKeySet)
	err := json.NewDecoder(reader).Decode(&jwks)
	if err != nil {
		return nil, err
	}

	set := &JSONWebKeySet{keys: map[string]*JSONWebKey{}}
	for _, k := range jwks.Keys {
		/*
		 * !Important!
		 * -----------
		 * copy k value onto local stack
		 * before k changes in the next iteration.
		 */
		k0 := k
		set.keys[k0.KeyID] = (*JSONWebKey)(&k0)
	}

	return set, nil
}

type JSONWebKeySet struct {
	keys map[string]*JSONWebKey
}

func (s *JSONWebKeySet) Size() int {
	if s == nil {
		return 0
	}
	return len(s.keys)
}

func (s *JSONWebKeySet) Merge(t *JSONWebKeySet) *JSONWebKeySet {
	if t == nil || t.Size() == 0 {
		return s
	}

	var keys []*JSONWebKey
	for _, k := range s.keys {
		keys = append(keys, k)
	}
	for _, k := range t.keys {
		keys = append(keys, k)
	}

	return NewJSONWebKeySet(keys...)
}

// FindByKeyID returns the JSONWebKey with matching kid header in the key set, or nil if no key should found.
func (s *JSONWebKeySet) FindByKeyID(kid string) *JSONWebKey {
	return s.keys[kid]
}

// FindForSigning returns the JSONWebKey intended for signature use with the matching algorithm, or nil if not found.
func (s *JSONWebKeySet) FindForSigning(alg spec.SignatureAlgorithm) *JSONWebKey {
	return s.find(alg.String(), UseSig)
}

// FindForEncryption returns the JSONWebKey intended for encryption use with the matching algorithm, or nil if not found.
func (s *JSONWebKeySet) FindForEncryption(alg spec.EncryptionAlgorithm) *JSONWebKey {
	return s.find(alg.String(), UseEnc)
}

func (s *JSONWebKeySet) find(alg string, use string) *JSONWebKey {
	var candidates []*JSONWebKey

	for _, k := range s.keys {
		if k.Use == use && k.Algorithm == alg {
			candidates = append(candidates, k)
		}
	}

	switch len(candidates) {
	case 0:
		return nil
	case 1:
		return candidates[0]
	default:
		return candidates[time.Now().Unix()%int64(len(candidates))]
	}
}

// Public returns a new JSONWebKeySet whose keys are stripped of all private or symmetric key information, suitable to
// be displayed publicly.
func (s *JSONWebKeySet) Public() *JSONWebKeySet {
	var pubKeys []*JSONWebKey
	for _, each := range s.keys {
		if pubKey := each.Public(); pubKey != nil {
			pubKeys = append(pubKeys, pubKey)
		}
	}
	return NewJSONWebKeySet(pubKeys...)
}

func (s *JSONWebKeySet) MarshalJSON() ([]byte, error) {
	jwks := &jose.JSONWebKeySet{Keys: []jose.JSONWebKey{}}
	for _, each := range s.keys {
		jwks.Keys = append(jwks.Keys, (jose.JSONWebKey)(*each))
	}
	return json.Marshal(jwks)
}

func (s *JSONWebKeySet) UnmarshalJSON(bytes []byte) error {
	jwks := new(jose.JSONWebKeySet)
	if err := json.Unmarshal(bytes, &jwks); err != nil {
		return err
	}

	if s.keys == nil {
		s.keys = map[string]*JSONWebKey{}
	}

	for _, each := range jwks.Keys {
		jwk := each // copy it onto the stack before referencing it, this should important!
		s.keys[jwk.KeyID] = (*JSONWebKey)(&jwk)
	}

	return nil
}
