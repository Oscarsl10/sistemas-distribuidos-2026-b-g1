package authsecurity

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	validEmail    = "client@example.test"
	validPassword = "correct horse"
)

type clock struct{ at time.Time }

func (c *clock) now() time.Time { return c.at }

// login sends POST /api/v1/auth/login as the client address the gateway appended to X-Forwarded-For.
func login(handler http.Handler, forwardedFor, password string) *httptest.ResponseRecorder {
	body := `{"email":"` + validEmail + `","password":"` + password + `"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	req.Header.Set("X-Forwarded-For", forwardedFor)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	return recorder
}

// Story S of the STRIDE analysis, criterion 1: the 11th login from one address is refused with 429
// even when the credentials are valid, and the answer tells when to retry.
func TestAuthzLoginRejectsTheEleventhFailedAttemptFromTheSameAddressEvenWithValidCredentials(t *testing.T) {
	handler := NewLoginHandler(NewLoginLimiter((&clock{time.Now()}).now), validEmail, validPassword)

	for i := 1; i <= 10; i++ {
		if got := login(handler, "198.51.100.7", "wrong").Code; got != http.StatusUnauthorized {
			t.Fatalf("attempt %d: got %d, want 401", i, got)
		}
	}
	answer := login(handler, "198.51.100.7", validPassword)

	if answer.Code != http.StatusTooManyRequests {
		t.Fatalf("11th attempt with valid credentials: got %d, want 429", answer.Code)
	}
	if answer.Header().Get("Retry-After") != "300" {
		t.Fatalf("Retry-After = %q, want 300", answer.Header().Get("Retry-After"))
	}
	if !strings.Contains(answer.Body.String(), "RATE_LIMITED") {
		t.Fatalf("body = %s", answer.Body.String())
	}
}

func TestAuthzLockOfOneAddressDoesNotLockAnotherAddress(t *testing.T) {
	handler := NewLoginHandler(NewLoginLimiter((&clock{time.Now()}).now), validEmail, validPassword)
	for i := 0; i < 10; i++ {
		login(handler, "198.51.100.7", "wrong")
	}
	if got := login(handler, "203.0.113.9", validPassword).Code; got != http.StatusOK {
		t.Fatalf("another address: got %d, want 200", got)
	}
}

func TestAuthzAddressWrittenByTheClientIsNotUsedToEscapeTheLock(t *testing.T) {
	handler := NewLoginHandler(NewLoginLimiter((&clock{time.Now()}).now), validEmail, validPassword)
	for i := 0; i < 10; i++ {
		login(handler, "198.51.100.7", "wrong")
	}
	// The client puts a fresh address first; the gateway appends the real one last.
	if got := login(handler, "192.0.2.1, 198.51.100.7", validPassword).Code; got != http.StatusTooManyRequests {
		t.Fatalf("spoofed first entry: got %d, want 429", got)
	}
}

func TestAuthzLockEndsAfterTheCooldown(t *testing.T) {
	now := &clock{time.Now()}
	handler := NewLoginHandler(NewLoginLimiter(now.now), validEmail, validPassword)
	for i := 0; i < 10; i++ {
		login(handler, "198.51.100.7", "wrong")
	}
	now.at = now.at.Add(5*time.Minute + time.Second)
	if got := login(handler, "198.51.100.7", validPassword).Code; got != http.StatusOK {
		t.Fatalf("after the cooldown: got %d, want 200", got)
	}
}
