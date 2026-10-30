package rehearsal

import (
	"crypto/rsa"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// The permission that csp-auth-api asks for on PATCH /users/{userId}/role.
const PermissionManageRoles = "user:role:update"

// NewService is the stand-in of csp-auth-api: it authenticates and authorizes by itself, whatever
// the caller says about who sent the request.
func NewService(key *rsa.PublicKey, now func() time.Time) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/auth/me", protect(key, now, "", func(w http.ResponseWriter, _ *http.Request, c Claims) {
		_, _ = w.Write([]byte(c.Subject))
	}))
	mux.HandleFunc("DELETE /api/v1/auth/users/{userId}", protect(key, now, PermissionManageRoles, func(w http.ResponseWriter, _ *http.Request, _ Claims) {
		w.WriteHeader(http.StatusNoContent)
	}))
	return mux
}

// protect answers 401 for a missing or invalid token and 403 for a valid token without the
// permission. The body is the same for every cause, so it tells an attacker nothing.
func protect(key *rsa.PublicKey, now func() time.Time, permission string, next func(http.ResponseWriter, *http.Request, Claims)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := Verify(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "), key, now())
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if permission != "" && !claims.Allows(permission) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r, claims)
	}
}

// NewEdge is the stand-in of csp-api-gateway: it drops the identity headers a client may send and
// only checks that Authorization is not empty, as nginx/routes/auth.conf does today.
func NewEdge(service http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("X-User-Id")
		r.Header.Del("X-User-Roles")
		if r.URL.Path != "/api/v1/auth/health" && r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		service.ServeHTTP(w, r)
	})
}

// HashPassword stores a password with bcrypt, the hash that csp-auth-api uses.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func PasswordMatches(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
