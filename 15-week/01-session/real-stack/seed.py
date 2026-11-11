"""Usage: W15_EMAIL=<synthetic email> W15_PASSWORD=<synthetic password> python seed.py
Creates one synthetic user in the copy of the release stack through the normal flow:
register -> link in Mailpit -> confirm -> login. Prints status codes only, never the credentials or the tokens."""
import json
import os
import re
import time
import urllib.request
import uuid

GATEWAY = "http://localhost:8000/api/v1/auth"
MAILPIT = "http://localhost:8025/api/v1"
EMAIL, PASSWORD = os.environ["W15_EMAIL"], os.environ["W15_PASSWORD"]  # a synthetic user of the copy; never a real account


def call(method, url, body=None, headers=None):
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(url, data=data, method=method, headers={"Content-Type": "application/json", **(headers or {})})
    try:
        with urllib.request.urlopen(req, timeout=15) as r:
            raw = r.read()
            return r.status, (json.loads(raw) if raw else None)
    except urllib.error.HTTPError as e:
        raw = e.read()
        try:
            return e.code, json.loads(raw) if raw else None
        except ValueError:
            return e.code, None


status, body = call("POST", GATEWAY + "/register", {"email": EMAIL, "password": PASSWORD, "name": "Carga Sintetica", "phone": "+570000000000", "address": "Calle 1 # 2-3"},
                    {"Idempotency-Key": str(uuid.uuid4())})
print("register ->", status)
token = None
for _ in range(20):
    time.sleep(1)
    _, listing = call("GET", MAILPIT + "/messages")
    for m in (listing or {}).get("messages", []):
        if any(EMAIL in (to.get("Address") or "") for to in m.get("To", [])):
            _, full = call("GET", MAILPIT + "/message/" + m["ID"])
            found = re.search(r"[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", (full or {}).get("Text", ""))
            if found:
                token = found.group(0)
                break
    if token:
        break
print("link in Mailpit ->", "found" if token else "NOT FOUND")
status, _ = call("POST", GATEWAY + "/register/confirm", {"token": token}, {"Idempotency-Key": str(uuid.uuid4())})
print("register/confirm ->", status)
status, body = call("POST", GATEWAY + "/login", {"email": EMAIL, "password": PASSWORD})
print("login ->", status, "| access token received:", bool(body and body.get("accessToken")))
