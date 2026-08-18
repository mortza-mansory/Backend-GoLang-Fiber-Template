// Package token provides JWT signing and verification infrastructure.
//
// This is infrastructure, not business logic. Modules that need tokens should
// depend on this package's narrow interface rather than on JWT specifics.
package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/yourorg/go-fiber-template/internal/config"
)

// Claims is the standard JWT claim set used by the template.
type Claims struct {
	Subject   string   `json:"sub"`
	Roles     []string `json:"roles,omitempty"`
	Issuer    string   `json:"iss,omitempty"`
	IssuedAt  int64    `json:"iat,omitempty"`
	ExpiresAt int64    `json:"exp,omitempty"`
}

// Valid checks the standard claims (exp) are acceptable.
func (c *Claims) Valid() error {
	if c.ExpiresAt > 0 && time.Now().Unix() > c.ExpiresAt {
		return errors.New("token: expired")
	}
	return nil
}

// GetExpirationTime implements jwt.Claims.
func (c *Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	if c.ExpiresAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(c.ExpiresAt, 0)), nil
}

// GetIssuedAt implements jwt.Claims.
func (c *Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	if c.IssuedAt == 0 {
		return nil, nil
	}
	return jwt.NewNumericDate(time.Unix(c.IssuedAt, 0)), nil
}

// GetIssuer implements jwt.Claims.
func (c *Claims) GetIssuer() (string, error) { return c.Issuer, nil }

// GetSubject implements jwt.Claims.
func (c *Claims) GetSubject() (string, error) { return c.Subject, nil }

// GetAudience implements jwt.Claims.
func (c *Claims) GetAudience() (jwt.ClaimStrings, error) { return nil, nil }

// NotBefore implements jwt.Claims.
func (c *Claims) GetNotBefore() (*jwt.NumericDate, error) { return nil, nil }

// Manager signs and verifies JWTs.
type Manager struct {
	secret []byte
	ttl    time.Duration
	issuer string
	method jwt.SigningMethod
}

// NewManager creates a JWT Manager from configuration.
func NewManager(cfg config.Security) (*Manager, error) {
	method := jwt.GetSigningMethod(cfg.JWTAlgorithm)
	if method == nil {
		return nil, fmt.Errorf("token: unsupported signing algorithm %q", cfg.JWTAlgorithm)
	}
	return &Manager{
		secret: []byte(cfg.JWTSecret),
		ttl:    cfg.JWTExpiry,
		issuer: cfg.JWTIssuer,
		method: method,
	}, nil
}

// Sign issues a signed JWT for the given subject and roles.
func (m *Manager) Sign(subject string, roles []string) (string, time.Time, error) {
	now := time.Now()
	exp := now.Add(m.ttl)
	claims := &Claims{
		Subject:   subject,
		Roles:     roles,
		Issuer:    m.issuer,
		IssuedAt:  now.Unix(),
		ExpiresAt: exp.Unix(),
	}
	token := jwt.NewWithClaims(m.method, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("token: sign: %w", err)
	}
	return signed, exp, nil
}

// Verify parses and validates a token string, returning its claims.
func (m *Manager) Verify(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != m.method {
			return nil, fmt.Errorf("token: unexpected signing method %v", t.Method)
		}
		return m.secret, nil
	}, jwt.WithIssuer(m.issuer))
	if err != nil {
		return nil, fmt.Errorf("token: verify: %w", err)
	}
	return claims, nil
}
