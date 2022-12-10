package spec

import (
	"encoding/json"
	"fmt"
)

const (
	DisplayPage Display = 1 << iota
	DisplayPopup
	DisplayWap
	DisplayTouch

	displayPage  = "page"
	displayPopup = "popup"
	displayWap   = "wap"
	displayTouch = "touch"
)

// Display represents the display parameter in OpenID Connect 1.0.
type Display uint8

func (d Display) String() string {
	switch d {
	case DisplayPage:
		return displayPage
	case DisplayPopup:
		return displayPopup
	case DisplayWap:
		return displayWap
	case DisplayTouch:
		return displayTouch
	default:
		return ""
	}
}

func (d Display) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Display) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch value {
	case displayPage:
		*d = DisplayPage
	case displayPopup:
		*d = DisplayPopup
	case displayWap:
		*d = DisplayWap
	case displayTouch:
		*d = DisplayTouch
	default:
		return fmt.Errorf("invalid spec.Display value [%s]", value)
	}

	return nil
}
