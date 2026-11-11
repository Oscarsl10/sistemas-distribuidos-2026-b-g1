package scale

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// Payments is the stand-in of the payment provider: a real HTTP server on a real port that
// can be killed and started again, with a bounded capacity and a fixed processing time.
type Payments struct {
	delay time.Duration
	slots chan struct{}
	seq   atomic.Int64

	mu   sync.Mutex
	addr string
	srv  *http.Server
}

// NewPayments returns a provider that takes delay to answer and serves at most capacity
// payments at the same time; the next ones wait for a free slot.
func NewPayments(delay time.Duration, capacity int) *Payments {
	return &Payments{delay: delay, slots: make(chan struct{}, capacity)}
}

// Start listens on the address of the previous start, or on a free port the first time.
func (p *Payments) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	addr := p.addr
	if addr == "" {
		addr = "127.0.0.1:0"
	}
	var ln net.Listener
	var err error
	for i := 0; i < 100; i++ {
		if ln, err = net.Listen("tcp", addr); err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if err != nil {
		return err
	}
	p.addr = ln.Addr().String()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /pay", p.pay)
	p.srv = &http.Server{Handler: mux}
	go func() { _ = p.srv.Serve(ln) }()
	return nil
}

// Kill stops the provider at once, as a crashed container does: the connections in progress are dropped.
func (p *Payments) Kill() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.srv != nil {
		_ = p.srv.Close()
		p.srv = nil
	}
}

func (p *Payments) URL() string { return "http://" + p.addr }

func (p *Payments) pay(w http.ResponseWriter, _ *http.Request) {
	p.slots <- struct{}{}
	defer func() { <-p.slots }()
	time.Sleep(p.delay)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"id": fmt.Sprintf("pay-%d", p.seq.Add(1)), "status": "APPROVED"})
}
