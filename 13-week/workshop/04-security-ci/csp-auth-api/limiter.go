// Package authsecurity is a minimal model of the login attempt limit of csp-auth-api (ADR-040):
// 10 failed attempts from one client address in 5 minutes lock that address for 5 minutes.
package authsecurity

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"
)

// LoginLimiter counts failed logins per client address. The real service keeps the counters in
// Redis (attempts:{ip} and locked:{ip}); this model keeps them in memory.
type LoginLimiter struct {
	mu          sync.Mutex
	now         func() time.Time
	maxAttempts int
	window      time.Duration
	lock        time.Duration
	failures    map[string][]time.Time
	lockedUntil map[string]time.Time
}

func NewLoginLimiter(now func() time.Time) *LoginLimiter {
	return &LoginLimiter{
		now: now, maxAttempts: 10, window: 5 * time.Minute, lock: 5 * time.Minute,
		failures: map[string][]time.Time{}, lockedUntil: map[string]time.Time{},
	}
}

// RetryAfter is the remaining seconds of the lock of the address, 0 when it is not locked.
func (l *LoginLimiter) RetryAfter(address string) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	remaining := l.lockedUntil[address].Sub(l.now())
	if remaining <= 0 {
		return 0
	}
	return int((remaining + time.Second - 1) / time.Second)
}

// Failed records one failed attempt and locks the address when it reaches the maximum.
func (l *LoginLimiter) Failed(address string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	recent := l.failures[address][:0]
	for _, at := range l.failures[address] {
		if now.Sub(at) < l.window {
			recent = append(recent, at)
		}
	}
	recent = append(recent, now)
	if len(recent) >= l.maxAttempts {
		l.lockedUntil[address] = now.Add(l.lock)
		recent = nil
	}
	l.failures[address] = recent
}

// ClientAddress is the last entry of X-Forwarded-For, the one the gateway appends (ADR-040 D8): the
// entries before it were written by the client and are never used.
func ClientAddress(r *http.Request) string {
	if lines := r.Header.Values("X-Forwarded-For"); len(lines) > 0 {
		entries := strings.Split(lines[len(lines)-1], ",")
		if address, err := netip.ParseAddr(strings.TrimSpace(entries[len(entries)-1])); err == nil {
			return address.Unmap().String()
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
