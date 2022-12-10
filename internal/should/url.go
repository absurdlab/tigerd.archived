package should

import (
	"errors"
	"net/url"
	"strings"
)

// URL returns a validation rule to check various aspects of the url. The rule accepts input types of string and
// *url.URL. It ignores any input type that is not compatible or that is nil or empty. By default, the rule does not
// allow any scheme (which means it required further configuration), allows fragment and query, and does not restrict
// host values.
func URL() *UrlRule {
	return &UrlRule{
		http:          false,
		https:         false,
		customScheme:  false,
		fragment:      true,
		query:         true,
		localhostOnly: false,
	}
}

type UrlRule struct {
	http          bool
	https         bool
	customScheme  bool
	fragment      bool
	query         bool
	localhostOnly bool
	message       string
}

// Http allows the url to have http scheme.
func (r *UrlRule) Http() *UrlRule {
	r.http = true
	return r
}

// Https allows the url to have https scheme.
func (r *UrlRule) Https() *UrlRule {
	r.https = true
	return r
}

// CustomScheme allows the url to use custom protocol, such as native application namespace.
func (r *UrlRule) CustomScheme() *UrlRule {
	r.customScheme = true
	return r
}

// NoFragment disallows the url to use fragment component.
func (r *UrlRule) NoFragment() *UrlRule {
	r.fragment = false
	return r
}

// NoQuery disallows the url to use query component.
func (r *UrlRule) NoQuery() *UrlRule {
	r.query = false
	return r
}

// LocalhostOnly restricts the url to have either "localhost" or "127.0.0.1" as hostname.
func (r *UrlRule) LocalhostOnly() *UrlRule {
	r.localhostOnly = true
	return r
}

// Error overrides default error messages. If set, any failure will return this message.
func (r *UrlRule) Error(message string) *UrlRule {
	r.message = message
	return r
}

func (r *UrlRule) Validate(value interface{}) error {
	switch v := value.(type) {
	case string:
		if len(v) == 0 {
			return nil
		}

		u, err := url.Parse(v)
		if err != nil {
			return err
		}

		return r.Validate(u)

	case *url.URL:
		if v == nil {
			return nil
		}

	default:
		return nil
	}

	u := value.(*url.URL)

	switch strings.ToLower(u.Scheme) {
	case "http":
		if !r.http {
			return r.error("should not use HTTP scheme")
		}
	case "https":
		if !r.https {
			return r.error("should not use HTTPS scheme")
		}
	default:
		if !r.customScheme {
			return r.error("should not use custom scheme")
		}
	}

	if !r.fragment && len(u.Fragment) > 0 {
		return r.error("should not have fragment component")
	}

	if !r.query && len(u.Query()) > 0 {
		return r.error("should not have query component")
	}

	if r.localhostOnly {
		switch strings.ToLower(u.Hostname()) {
		case "localhost", "127.0.0.1":
		default:
			return r.error("should have localhost as host")
		}
	}

	return nil
}

func (r *UrlRule) error(message string) error {
	if len(r.message) > 0 {
		return errors.New(r.message)
	}
	return errors.New(message)
}
