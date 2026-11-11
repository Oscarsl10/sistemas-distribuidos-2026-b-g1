"""Tests of the week 15 session 1 against the COPY of the release stack (project csp-w15, fresh data).
Usage: W15_EMAIL=... W15_PASSWORD=... python tests.py <gateway-credentials | gateway-perclient | capacity | redis-down | kill-auth>
Prints status codes and times only, never the credentials or the tokens."""
import json
import os
import subprocess
import sys
import threading
import time
import urllib.error
import urllib.request
from collections import Counter
from concurrent.futures import ThreadPoolExecutor

GW = "http://localhost:8000/api/v1/auth"
DIRECT = "http://127.0.0.1:18081/api/v1/auth"
EMAIL, PASSWORD = os.environ["W15_EMAIL"], os.environ["W15_PASSWORD"]  # a synthetic user of the copy; never a real account


def req(method, url, body=None, timeout=10):
    t = time.perf_counter()
    data = json.dumps(body).encode() if body is not None else None
    r = urllib.request.Request(url, data=data, method=method, headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(r, timeout=timeout) as resp:
            resp.read()
            return resp.status, (time.perf_counter() - t) * 1000, dict(resp.headers)
    except urllib.error.HTTPError as e:
        e.read()
        return e.code, (time.perf_counter() - t) * 1000, dict(e.headers)
    except Exception:
        return 0, (time.perf_counter() - t) * 1000, {}


def good_login(base):
    return lambda: req("POST", base + "/login", {"email": EMAIL, "password": PASSWORD})


def bad_login(base):
    return lambda: req("POST", base + "/login", {"email": EMAIL, "password": "Wrong-Pass-1"})


def jwks(base):
    return lambda: req("GET", base + "/jwks")


def open_loop(fn, rate, dur, hook=None):
    out, lock, t0 = [], threading.Lock(), time.perf_counter()

    def one(rel):
        s, l, _ = fn()
        with lock:
            out.append((rel, s, l))

    with ThreadPoolExecutor(max_workers=300) as ex:
        i = 0
        while i / rate < dur:
            due = t0 + i / rate
            now = time.perf_counter()
            if due > now:
                time.sleep(due - now)
            rel = time.perf_counter() - t0
            if hook:
                hook(rel)
            ex.submit(one, rel)
            i += 1
    return out


def pct(values, p):
    v = sorted(values)
    return v[int((len(v) - 1) * p)] if v else 0


def docker(*args):
    return subprocess.run(["docker", *args], capture_output=True, text=True).stdout.strip()


def sequential(fn, n):
    rs = [fn() for _ in range(n)]
    return [s for s, _, _ in rs], rs


def gateway_credentials():
    print("== A. the credential bucket of the gateway (2 r/s, burst 3) on POST /login, valid credentials")
    statuses, rs = sequential(good_login(GW), 12)
    retry = next((h.get("Retry-After") for s, _, h in rs if s == 429), None)
    print("12 requests in a row:", " ".join(map(str, statuses)), "| Retry-After of a 429:", retry)
    time.sleep(4)
    out = open_loop(good_login(GW), 5, 10)
    c = Counter(s for _, s, _ in out)
    print("sustained 5 req/s for 10 s: sent", len(out), "| answers", dict(sorted(c.items())), "| expected about 4 + 2 x 10 = 24 accepted")


def gateway_perclient():
    print("== B. the per-client limit of the gateway (10 r/s, burst 20) on GET /jwks")
    print(f"{'rate':>6} {'sent':>5} {'200':>5} {'429':>5} {'other':>5} {'p95 ms':>8}")
    for rate in (8, 15, 30):
        time.sleep(3)
        out = open_loop(jwks(GW), rate, 6)
        c = Counter(s for _, s, _ in out)
        other = len(out) - c.get(200, 0) - c.get(429, 0)
        print(f"{rate:>6} {len(out):>5} {c.get(200, 0):>5} {c.get(429, 0):>5} {other:>5} {pct([l for _, _, l in out], .95):>8.0f}")


def capacity():
    print("== C. capacity of csp-auth-api, direct (no gateway), POST /login with valid credentials; SLO of ADR-061: p95 <= 1000 ms")
    print(f"{'rate':>5} {'sent':>5} {'200':>5} {'429':>4} {'503':>4} {'other':>5} {'p50':>6} {'p95':>6} {'p99':>6} {'p95 of 200':>10}  SLO")
    for rate in (5, 10, 15, 20, 30, 40):
        time.sleep(3)
        out = open_loop(good_login(DIRECT), rate, 10)
        c = Counter(s for _, s, _ in out)
        lat = [l for _, _, l in out]
        ok = [l for _, s, l in out if s == 200]
        other = len(out) - c.get(200, 0) - c.get(429, 0) - c.get(503, 0)
        met = pct(lat, .95) <= 1000 and c.get(200, 0) == len(out)
        print(f"{rate:>5} {len(out):>5} {c.get(200, 0):>5} {c.get(429, 0):>4} {c.get(503, 0):>4} {other:>5} {pct(lat, .5):>6.0f} {pct(lat, .95):>6.0f} {pct(lat, .99):>6.0f} {pct(ok, .95):>10.0f}  met={met}")


def redis_down():
    print("== D. chaos: Redis stopped while the logins go on (direct to csp-auth-api)")
    docker("stop", "csp-w15-redis-1")
    print("redis stopped:", docker("ps", "-a", "--filter", "name=csp-w15-redis-1", "--format", "{{.Status}}"))
    t = time.perf_counter()
    statuses, rs = sequential(bad_login(DIRECT), 15)
    print("15 wrong-password logins with Redis down:", " ".join(map(str, statuses)), "| total", f"{time.perf_counter() - t:.1f}s",
          "| first request", f"{rs[0][1]:.0f} ms, the rest p95 {pct([l for _, l, _ in rs[1:]], .95):.0f} ms")
    s2, rs2 = sequential(good_login(DIRECT), 3)
    print("3 valid logins with Redis down:", " ".join(map(str, s2)))
    docker("start", "csp-w15-redis-1")
    for _ in range(30):
        if "healthy" in docker("ps", "--filter", "name=csp-w15-redis-1", "--format", "{{.Status}}"):
            break
        time.sleep(1)
    print("redis back:", docker("ps", "--filter", "name=csp-w15-redis-1", "--format", "{{.Status}}"), "| waiting 12 s for the pause of the breaker (10 s) to end")
    time.sleep(12)
    s3, _ = sequential(bad_login(DIRECT), 13)
    print("control, Redis up again, 13 wrong-password logins:", " ".join(map(str, s3)), "(the lock of the 11th attempt is the criterion of story #280)")
    logs = docker("logs", "--since", "2m", "csp-w15-auth-api-1")
    warns = [json.loads(x) for x in logs.splitlines() if x.startswith("{") and '"level":"WARN"' in x]
    print("WARN lines of auth-api in the last 2 minutes:", len(warns), "| messages:", sorted({w.get("msg") for w in warns}))


def kill_auth():
    print("== E. chaos: csp-auth-api killed in the middle of the traffic through the gateway (GET /jwks at 8 req/s)")
    state = {"killed": None, "started": None}

    def hook(rel):
        if state["killed"] is None and rel >= 6:
            state["killed"] = rel
            docker("kill", "csp-w15-auth-api-1")
        if state["killed"] is not None and state["started"] is None and rel >= 12:
            state["started"] = rel
            print("  at t=%.1f s: docker start (restart policy: %s)" % (rel, docker("inspect", "csp-w15-auth-api-1", "--format", "{{.HostConfig.RestartPolicy.Name}}")),
                  "| state before:", docker("inspect", "csp-w15-auth-api-1", "--format", "{{.State.Status}}"))
            docker("start", "csp-w15-auth-api-1")

    out = sorted(open_loop(jwks(GW), 8, 26, hook))
    print(f"{'window':>12} {'sent':>5} {'200':>5} {'502':>5} {'504':>5} {'other':>5} {'p95 ms':>7}")
    for a, b, name in ((0, 6, "before kill"), (6, 12, "killed"), (12, 18, "after start"), (18, 26, "recovered")):
        w = [(s, l) for rel, s, l in out if a <= rel < b]
        c = Counter(s for s, _ in w)
        print(f"{name:>12} {len(w):>5} {c.get(200, 0):>5} {c.get(502, 0):>5} {c.get(504, 0):>5} {len(w) - c.get(200, 0) - c.get(502, 0) - c.get(504, 0):>5} {pct([l for _, l in w], .95):>7.0f}")
    after = [rel for rel, s, _ in out if rel >= state["started"] and s == 200]
    print(f"first 200 after docker start: {after[0] - state['started']:.1f} s" if after else "no 200 after the start")
    fails = [l for rel, s, l in out if s != 200 and state["killed"] <= rel < state["started"]]
    print(f"answers during the outage: {len(fails)}, p95 {pct(fails, .95):.0f} ms, max {max(fails) if fails else 0:.0f} ms")


if __name__ == "__main__":
    {"gateway-credentials": gateway_credentials, "gateway-perclient": gateway_perclient, "capacity": capacity,
     "redis-down": redis_down, "kill-auth": kill_auth}[sys.argv[1]]()
