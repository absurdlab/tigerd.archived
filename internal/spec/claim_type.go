package spec

import (
	"encoding/json"
	"fmt"
)

const (
	ClaimTypeNormal ClaimType = 1 << iota
	ClaimTypeAggregated
	ClaimTypeDistributed

	claimTypeNormal      = "normal"
	claimTypeAggregated  = "aggregated"
	claimTypeDistributed = "distributed"
)

// ClaimType represents the claim_type parameter.
type ClaimType uint8

func (t ClaimType) String() string {
	switch t {
	case ClaimTypeNormal:
		return claimTypeNormal
	case ClaimTypeAggregated:
		return claimTypeAggregated
	case ClaimTypeDistributed:
		return claimTypeDistributed
	default:
		return ""
	}
}

func (t ClaimType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

func (t *ClaimType) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case claimTypeNormal:
		*t = ClaimTypeNormal
	case claimTypeAggregated:
		*t = ClaimTypeAggregated
	case claimTypeDistributed:
		*t = ClaimTypeDistributed
	default:
		return fmt.Errorf("invalid value for spec.ClaimType [%s]", value)
	}

	return nil
}
