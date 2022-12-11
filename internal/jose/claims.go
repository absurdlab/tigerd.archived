package jose

import (
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/oklog/ulid/v2"
	"time"
)

// MultipleClaims can be implemented by a single structure to provide more than one claims sources. All sources will be
// aggregated and encoded as a flat structure in the final JWT payload. Name collisions are not handled.
type MultipleClaims interface {
	MultipleClaims() []any
}

// StdClaims aliases jwt.Claims to provide a fluent builder API.
//
//		// Example:
//		new(StdClaims).
//			GenerateID().
//			WithIssuer("foo").
//			WithIssuedAtNow().
//			WithSubject("bar").
//			WithAudience("one", "two", "three").
//			WithExpiryIn(2*time.Hour)
type StdClaims jwt.Claims

func (c *StdClaims) GenerateID() *StdClaims {
	c.ID = ulid.Make().String()
	return c
}

func (c *StdClaims) WithID(jti string) *StdClaims {
	c.ID = jti
	return c
}

func (c *StdClaims) WithIssuer(iss string) *StdClaims {
	c.Issuer = iss
	return c
}

func (c *StdClaims) WithIssuedAt(iat time.Time) *StdClaims {
	c.IssuedAt = jwt.NewNumericDate(iat)
	return c
}

func (c *StdClaims) WithIssuedAtNow() *StdClaims {
	return c.WithIssuedAt(time.Now())
}

func (c *StdClaims) WithSubject(sub string) *StdClaims {
	c.Subject = sub
	return c
}

func (c *StdClaims) WithAudience(aud ...string) *StdClaims {
	if c.Audience == nil {
		c.Audience = []string{}
	}
	if len(aud) > 0 {
		c.Audience = append(c.Audience, aud...)
	}
	return c
}

func (c *StdClaims) WithExpiry(exp time.Time) *StdClaims {
	c.Expiry = jwt.NewNumericDate(exp)
	return c
}

func (c *StdClaims) WithExpiryIn(dur time.Duration) *StdClaims {
	return c.WithExpiry(time.Now().Add(dur))
}

func (c *StdClaims) WithNotBefore(nbf time.Time) *StdClaims {
	c.NotBefore = jwt.NewNumericDate(nbf)
	return c
}

func (c *StdClaims) WithNotBeforeIn(dur time.Duration) *StdClaims {
	return c.WithNotBefore(time.Now().Add(dur))
}
