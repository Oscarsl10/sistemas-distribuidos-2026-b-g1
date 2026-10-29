package authsecurity

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// LoginRoute is the path of the login through the gateway.
const LoginRoute = "POST /api/v1/auth/login"

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewLoginHandler models POST /api/v1/auth/login: the lock is checked before the password, so a
// locked address is rejected with 429 RATE_LIMITED and Retry-After even with valid credentials.
func NewLoginHandler(limiter *LoginLimiter, validEmail, validPassword string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		address := ClientAddress(r)
		if seconds := limiter.RetryAfter(address); seconds > 0 {
			w.Header().Set("Retry-After", strconv.Itoa(seconds))
			writeError(w, http.StatusTooManyRequests, "RATE_LIMITED")
			return
		}
		var body loginRequest
		if json.NewDecoder(r.Body).Decode(&body) != nil || body.Email != validEmail || body.Password != validPassword {
			limiter.Failed(address)
			writeError(w, http.StatusUnauthorized, "INVALID_CREDENTIALS")
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

func writeError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": code})
}
