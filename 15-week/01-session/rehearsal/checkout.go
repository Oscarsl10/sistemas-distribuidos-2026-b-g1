package scale

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"
)

const (
	Confirmed   = "CONFIRMED"   // reserved, paid and confirmed
	Compensated = "COMPENSATED" // the payment failed and the reservation was released
	Rejected    = "REJECTED"    // the breaker is open: nothing was reserved
)

// Order is the answer of the checkout.
type Order struct {
	Status string `json:"status"`
}

// Breaker opens after threshold failures in a row, refuses calls for cooldown, and then lets a
// single probe through: a success closes it, a failure keeps it open for another cooldown.
type Breaker struct {
	threshold int
	cooldown  time.Duration

	mu       sync.Mutex
	failures int
	openedAt time.Time
	probing  bool
}

func NewBreaker(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{threshold: threshold, cooldown: cooldown}
}

func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.failures < b.threshold {
		return true
	}
	if !b.probing && time.Since(b.openedAt) >= b.cooldown {
		b.probing = true
		return true
	}
	return false
}

func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures, b.probing = 0, false
}

func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold {
		b.openedAt, b.probing = time.Now(), false
	}
}

// Checkout is the money path of the rehearsal: reserve a seat, pay, confirm; if the payment
// fails, release the seat (the compensation of a saga).
type Checkout struct {
	payments string
	client   *http.Client
	breaker  *Breaker

	mu        sync.Mutex
	stock     int
	reserved  int
	confirmed int
}

func NewCheckout(paymentsURL string, stock int, breaker *Breaker) *Checkout {
	return &Checkout{
		payments: paymentsURL,
		client:   &http.Client{Timeout: time.Second},
		breaker:  breaker,
		stock:    stock,
	}
}

// Counters returns the seats that are free, held and confirmed.
func (c *Checkout) Counters() (free, held, confirmed int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stock - c.reserved - c.confirmed, c.reserved, c.confirmed
}

func (c *Checkout) PlaceOrder(ctx context.Context) Order {
	if !c.breaker.Allow() {
		return Order{Status: Rejected}
	}
	c.mu.Lock()
	c.reserved++
	c.mu.Unlock()
	if err := c.pay(ctx); err != nil {
		c.breaker.Failure()
		c.mu.Lock()
		c.reserved--
		c.mu.Unlock()
		return Order{Status: Compensated}
	}
	c.breaker.Success()
	c.mu.Lock()
	c.reserved--
	c.confirmed++
	c.mu.Unlock()
	return Order{Status: Confirmed}
}

func (c *Checkout) pay(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.payments+"/pay", nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("payment refused")
	}
	return nil
}

// Handler exposes POST /orders: 201 confirmed, 502 compensated, 503 rejected by the open breaker.
func (c *Checkout) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		order := c.PlaceOrder(r.Context())
		status := map[string]int{Confirmed: http.StatusCreated, Compensated: http.StatusBadGateway, Rejected: http.StatusServiceUnavailable}[order.Status]
		if order.Status == Rejected {
			w.Header().Set("Retry-After", "1")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(order)
	})
	return mux
}
