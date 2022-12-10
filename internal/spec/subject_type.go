package spec

import (
	"encoding/json"
	"fmt"
)

const (
	SubjectTypePublic SubjectType = 1 << iota
	SubjectTypePairwise

	subjectTypePublic   = "public"
	subjectTypePairwise = "pairwise"
)

// SubjectType represents the subject_type parameter in OpenID Connect 1.0.
type SubjectType uint8

func (t SubjectType) String() string {
	switch t {
	case SubjectTypePublic:
		return subjectTypePublic
	case SubjectTypePairwise:
		return subjectTypePairwise
	default:
		return ""
	}
}

func (t SubjectType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *SubjectType) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case subjectTypePublic:
		*t = SubjectTypePublic
	case subjectTypePairwise:
		*t = SubjectTypePairwise
	default:
		return fmt.Errorf("invalid value for spec.SubjectType [%s]", value)
	}

	return nil
}
