package rehearsal

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var (
	now     = time.Date(2026, 10, 30, 10, 0, 0, 0, time.UTC)
	signer  = mustKey()
	service = NewService(&signer.PublicKey, func() time.Time { return now })
	edge    = NewEdge(service)
)

func mustKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return key
}

func b64(v any) string {
	raw, _ := json.Marshal(v)
	return base64.RawURLEncoding.EncodeToString(raw)
}

// token signs claims with the given key and the algorithm name written in the header.
func token(key *rsa.PrivateKey, alg string, c Claims) string {
	body := b64(map[string]string{"alg": alg, "typ": "JWT"}) + "." + b64(c)
	digest := sha256.Sum256([]byte(body))
	sig, _ := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	return body + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func claims(permissions ...string) Claims {
	return Claims{Subject: "user-1", Issuer: Issuer, Audience: Audience, ExpiresAt: now.Add(15 * time.Minute).Unix(), Permissions: permissions}
}

func call(h http.Handler, method, path, bearer string) int {
	req := httptest.NewRequest(method, path, nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code
}

const adminRoute = "/api/v1/auth/users/u-9"

func TestAuthzNoTokenAnswers401(t *testing.T) {
	if got := call(service, "DELETE", adminRoute, ""); got != 401 {
		t.Fatalf("got %d", got)
	}
}

func TestAuthzValidTokenWithoutThePermissionAnswers403(t *testing.T) {
	if got := call(service, "DELETE", adminRoute, token(signer, "RS256", claims("profile:read"))); got != 403 {
		t.Fatalf("got %d", got)
	}
}

func TestAuthzTokenWithThePermissionIsAccepted(t *testing.T) {
	if got := call(service, "DELETE", adminRoute, token(signer, "RS256", claims(PermissionManageRoles))); got != 204 {
		t.Fatalf("got %d", got)
	}
}

func TestAuthzTokenSignedByAnotherKeyAnswers401(t *testing.T) {
	if got := call(service, "GET", "/api/v1/auth/me", token(mustKey(), "RS256", claims())); got != 401 {
		t.Fatalf("got %d", got)
	}
}

func TestAuthzTokenThatNamesAnotherAlgorithmAnswers401(t *testing.T) {
	for _, alg := range []string{"none", "HS256"} {
		if got := call(service, "GET", "/api/v1/auth/me", token(signer, alg, claims())); got != 401 {
			t.Fatalf("alg %s: got %d", alg, got)
		}
	}
}

func TestAuthzExpiredTokenAnswers401(t *testing.T) {
	c := claims()
	c.ExpiresAt = now.Add(-time.Second).Unix()
	if got := call(service, "GET", "/api/v1/auth/me", token(signer, "RS256", c)); got != 401 {
		t.Fatalf("got %d", got)
	}
}

func TestAuthzTokenForAnotherAudienceAnswers401(t *testing.T) {
	c := claims()
	c.Audience = "catalog-service"
	if got := call(service, "GET", "/api/v1/auth/me", token(signer, "RS256", c)); got != 401 {
		t.Fatalf("got %d", got)
	}
}

// The edge only checks that the header is there: a made-up token passes it and the service stops it.
func TestAuthzTheEdgeAloneDoesNotValidateTheToken(t *testing.T) {
	made := "Bearer made.up.token"
	req := httptest.NewRequest("DELETE", adminRoute, nil)
	req.Header.Set("Authorization", made)
	rec := httptest.NewRecorder()
	edge.ServeHTTP(rec, req)
	if rec.Code != 401 || strings.Contains(rec.Body.String(), "made") {
		t.Fatalf("got %d %q", rec.Code, rec.Body.String())
	}
}

// Zero trust: a caller that skips the edge (a compromised service, an SSRF) meets the same checks.
func TestAuthzCallingTheServiceDirectlyIsStillRejected(t *testing.T) {
	req := httptest.NewRequest("DELETE", adminRoute, nil)
	req.Header.Set("X-User-Roles", "ADMIN")
	rec := httptest.NewRecorder()
	service.ServeHTTP(rec, req)
	if rec.Code != 401 {
		t.Fatalf("got %d", rec.Code)
	}
}

func TestAuthzTheEdgeDropsTheIdentityHeadersOfTheClient(t *testing.T) {
	var seen string
	probe := NewEdge(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) { seen = r.Header.Get("X-User-Roles") }))
	req := httptest.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer x")
	req.Header.Set("X-User-Roles", "ADMIN")
	probe.ServeHTTP(httptest.NewRecorder(), req)
	if seen != "" {
		t.Fatalf("the service received %q", seen)
	}
}

func TestPasswordIsStoredHashedAndStillVerifies(t *testing.T) {
	hash, err := HashPassword("S3cret-passphrase")
	if err != nil || strings.Contains(hash, "S3cret") || !PasswordMatches(hash, "S3cret-passphrase") || PasswordMatches(hash, "other") {
		t.Fatalf("hash %q err %v", hash, err)
	}
}
