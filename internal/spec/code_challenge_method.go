package spec

import (
	"encoding/json"
	"fmt"
)

const (
	CodeChallengeMethodPlain CodeChallengeMethod = 1 << iota
	CodeChallengeMethodS256

	codeChallengeMethodPlain = "plain"
	codeChallengeMethodS256  = "S256"
)

// CodeChallengeMethod represents the code_challenge_method parameter used in PKCE.
type CodeChallengeMethod uint8

func (m CodeChallengeMethod) String() string {
	switch m {
	case CodeChallengeMethodPlain:
		return codeChallengeMethodPlain
	case CodeChallengeMethodS256:
		return codeChallengeMethodS256
	default:
		return ""
	}
}

func (m CodeChallengeMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.String())
}

func (m *CodeChallengeMethod) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case codeChallengeMethodPlain:
		*m = CodeChallengeMethodPlain
	case codeChallengeMethodS256:
		*m = CodeChallengeMethodS256
	default:
		return fmt.Errorf("invalid value for spec.CodeChallengeMethod [%s]", value)
	}

	return nil
}
