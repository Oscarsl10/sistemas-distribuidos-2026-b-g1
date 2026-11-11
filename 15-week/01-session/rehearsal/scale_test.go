package scale

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"sync"
	"testing"
	"time"
)

const stock = 1_000_000

func startPayments(t *testing.T, delay time.Duration, capacity int) *Payments {
	t.Helper()
	p := NewPayments(delay, capacity)
	if err := p.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(p.Kill)
	return p
}

func newClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{MaxIdleConnsPerHost: 2000}}
}

// call places one order through the HTTP entry and returns the status code and the time it took.
func call(client *http.Client, url string) (int, time.Duration) {
	start := time.Now()
	resp, err := client.Post(url+"/orders", "application/json", nil)
	if err != nil {
		return 0, time.Since(start)
	}
	resp.Body.Close()
	return resp.StatusCode, time.Since(start)
}

func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	return sorted[int(float64(len(sorted)-1)*p)]
}

// result is one request of a run: when it started, what it answered and how long it took.
type result struct {
	at     time.Duration
	status int
	took   time.Duration
}

// run sends requests at a fixed rate, whatever the speed of the answers (open loop), and waits for all of them.
func run(client *http.Client, url string, rate int, d time.Duration, during func(elapsed time.Duration)) []result {
	var mu sync.Mutex
	var wg sync.WaitGroup
	var out []result
	t0 := time.Now()
	tick := time.NewTicker(time.Second / time.Duration(rate))
	defer tick.Stop()
	for time.Since(t0) < d {
		<-tick.C
		at := time.Since(t0)
		if during != nil {
			during(at)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			status, took := call(client, url)
			mu.Lock()
			out = append(out, result{at, status, took})
			mu.Unlock()
		}()
	}
	wg.Wait()
	return out
}

func summarize(rs []result) (n, ok int, p50, p95, p99 time.Duration) {
	var lat []time.Duration
	for _, r := range rs {
		lat = append(lat, r.took)
		if r.status == http.StatusCreated {
			ok++
		}
	}
	slices.Sort(lat)
	return len(rs), ok, percentile(lat, .5), percentile(lat, .95), percentile(lat, .99)
}

// Contract: the checkout (consumer) expects this answer from the payment provider.
func TestContractPaymentsAnswersWhatCheckoutExpects(t *testing.T) {
	raw, err := os.ReadFile("contract/payments.contract.json")
	if err != nil {
		t.Fatal(err)
	}
	var c struct {
		Request  struct{ Method, Path string }
		Response struct {
			Status         int
			ContentType    string
			RequiredFields []string
			AllowedStatus  []string
		}
	}
	if err := json.Unmarshal(raw, &c); err != nil {
		t.Fatal(err)
	}
	p := startPayments(t, time.Millisecond, 4)
	req, _ := http.NewRequest(c.Request.Method, p.URL()+c.Request.Path, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if resp.StatusCode != c.Response.Status || resp.Header.Get("Content-Type") != c.Response.ContentType {
		t.Fatalf("status %d, content type %q", resp.StatusCode, resp.Header.Get("Content-Type"))
	}
	for _, f := range c.Response.RequiredFields {
		if body[f] == "" {
			t.Fatalf("the field %q is missing in %v", f, body)
		}
	}
	if !slices.Contains(c.Response.AllowedStatus, body["status"]) {
		t.Fatalf("status %q is not in the contract", body["status"])
	}
	t.Logf("contract met: %d %s, fields %v, status %s", resp.StatusCode, resp.Header.Get("Content-Type"), c.Response.RequiredFields, body["status"])
}

// E2E: the money path through the real HTTP entry and the real provider: reserve, pay, confirm.
func TestE2EPlaceOrderReservePayConfirm(t *testing.T) {
	p := startPayments(t, time.Millisecond, 4)
	c := NewCheckout(p.URL(), stock, NewBreaker(5, time.Second))
	srv := httptest.NewServer(c.Handler())
	defer srv.Close()
	status, _ := call(newClient(), srv.URL)
	free, held, confirmed := c.Counters()
	if status != http.StatusCreated || confirmed != 1 || held != 0 || free != stock-1 {
		t.Fatalf("status %d, free %d, held %d, confirmed %d", status, free, held, confirmed)
	}
	t.Logf("order confirmed: status %d, free %d, held %d, confirmed %d", status, free, held, confirmed)
}

// Load: raise the rate until the latency leaves the SLO (p95 at most 250 ms, no failed order) and report where.
func TestLoadFindsTheBreakingPointAgainstTheSLO(t *testing.T) {
	const sloP95, target = 250 * time.Millisecond, 200
	p := startPayments(t, 50*time.Millisecond, 20) // capacity of the provider: 20 / 50 ms = 400 payments per second
	client := newClient()
	t.Logf("%-10s %6s %6s %9s %9s %9s  %s", "rate", "sent", "ok", "p50", "p95", "p99", "SLO")
	breaking := 0
	for _, rate := range []int{100, 200, 300, 400, 600} {
		c := NewCheckout(p.URL(), stock, NewBreaker(5, time.Second))
		srv := httptest.NewServer(c.Handler())
		n, ok, p50, p95, p99 := summarize(run(client, srv.URL, rate, 2*time.Second, nil))
		srv.Close()
		met := p95 <= sloP95 && ok == n
		if !met && breaking == 0 {
			breaking = rate
		}
		if rate == target && !met {
			t.Fatalf("at the target of %d req/s: ok %d of %d, p95 %v (SLO %v)", rate, ok, n, p95, sloP95)
		}
		t.Logf("%-10d %6d %6d %9v %9v %9v  met=%v", rate, n, ok, p50.Round(time.Millisecond), p95.Round(time.Millisecond), p99.Round(time.Millisecond), met)
	}
	t.Logf("SLO: p95 <= %v and every order confirmed. Target %d req/s met; first rate that breaks it: %d req/s", sloP95, target, breaking)
}

// Chaos: kill the provider in the middle of the traffic, bring it back, and check the system degrades and recovers.
func TestChaosKillThePaymentsServiceMidFlow(t *testing.T) {
	p := startPayments(t, 50*time.Millisecond, 20)
	c := NewCheckout(p.URL(), stock, NewBreaker(5, 500*time.Millisecond))
	srv := httptest.NewServer(c.Handler())
	defer srv.Close()
	killed, restarted := false, false
	rs := run(newClient(), srv.URL, 100, 5*time.Second, func(at time.Duration) {
		if !killed && at >= time.Second {
			killed = true
			p.Kill()
		}
		if killed && !restarted && at >= 3*time.Second {
			restarted = true
			if err := p.Start(); err != nil {
				t.Error(err)
			}
		}
	})
	phase := func(from, to time.Duration) (confirmed, compensated, rejected int, rejectedP95 time.Duration) {
		var lat []time.Duration
		for _, r := range rs {
			if r.at < from || r.at >= to {
				continue
			}
			switch r.status {
			case http.StatusCreated:
				confirmed++
			case http.StatusBadGateway:
				compensated++
			case http.StatusServiceUnavailable:
				rejected++
				lat = append(lat, r.took)
			}
		}
		slices.Sort(lat)
		return confirmed, compensated, rejected, percentile(lat, .95)
	}
	t.Logf("%-22s %10s %12s %9s %12s", "phase", "confirmed", "compensated", "rejected", "rejected p95")
	for _, ph := range []struct {
		name     string
		from, to time.Duration
	}{{"before the kill", 0, time.Second}, {"provider down", time.Second, 3 * time.Second}, {"provider back", 3 * time.Second, 4 * time.Second}, {"recovered", 4 * time.Second, 6 * time.Second}} {
		a, b, r, p95 := phase(ph.from, ph.to)
		t.Logf("%-22s %10d %12d %9d %12v", ph.name, a, b, r, p95.Round(time.Millisecond))
	}
	// The orders in flight when the provider dies (about 100 req/s x 50 ms = 5) are compensated; nothing else fails before the kill.
	if a, b, r, _ := phase(0, time.Second); b > 10 || r != 0 || a == 0 {
		t.Errorf("before the kill: confirmed %d, compensated %d (at most 10, the orders in flight), rejected %d", a, b, r)
	}
	_, comp, rej, rejP95 := phase(time.Second, 3*time.Second)
	if comp > 30 || rej < 100 || rejP95 > 50*time.Millisecond {
		t.Errorf("during the outage: compensated %d (at most 30), rejected %d (at least 100), rejected p95 %v (at most 50ms)", comp, rej, rejP95)
	}
	if _, b, r, _ := phase(4*time.Second, 6*time.Second); b != 0 || r != 0 {
		t.Errorf("after recovery: compensated %d, rejected %d", b, r)
	}
	free, held, confirmed := c.Counters()
	total := 0
	for _, r := range rs {
		if r.status == http.StatusCreated {
			total++
		}
	}
	if held != 0 || confirmed != total || free != stock-confirmed {
		t.Errorf("seats: free %d, held %d, confirmed %d, orders answered 201: %d", free, held, confirmed, total)
	}
	t.Logf("no seat is left held: free %d, held %d, confirmed %d (= orders answered 201: %d)", free, held, confirmed, total)
}
