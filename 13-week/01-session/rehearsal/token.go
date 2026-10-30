package rehearsal

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// The registered claims of csp-auth-api (issuerName and audience of its token package).
const (
	Issuer   = "csp-auth"
	Audience = "csp-api"
)

var ErrInvalidToken = errors.New("invalid token")

// Claims is the part of the AccessTokenClaims of auth-service.yaml that the rehearsal needs.
type Claims struct {
	Subject     string   `json:"sub"`
	Issuer      string   `json:"iss"`
	Audience    string   `json:"aud"`
	ExpiresAt   int64    `json:"exp"`
	Permissions []string `json:"permissions"`
}

func (c Claims) Allows(permission string) bool {
	for _, p := range c.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// Verify checks a token the way a service must, without asking anybody: the algorithm is fixed
// to RS256 (never the one the token asks for), then the signature, the issuer, the audience and
// the expiry. It never calls the database.
func Verify(token string, key *rsa.PublicKey, now time.Time) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}
	var header struct {
		Alg string `json:"alg"`
	}
	if decodeJSON(parts[0], &header) != nil || header.Alg != "RS256" {
		return Claims{}, ErrInvalidToken
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	digest := sha256.Sum256([]byte(parts[0] + "." + parts[1]))
	if rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], signature) != nil {
		return Claims{}, ErrInvalidToken
	}
	var claims Claims
	if decodeJSON(parts[1], &claims) != nil {
		return Claims{}, ErrInvalidToken
	}
	if claims.Issuer != Issuer || claims.Audience != Audience || !now.Before(time.Unix(claims.ExpiresAt, 0)) {
		return Claims{}, ErrInvalidToken
	}
	return claims, nil
}

func decodeJSON(segment string, into any) error {
	raw, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, into)
}
