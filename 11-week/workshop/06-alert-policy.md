# CP7 - Alert policy

Rules: [06-alerts.yml](06-alerts.yml), checked with `promtool check rules` (2 rules found). They alert on the availability SLO
of the registration (`POST /api/v1/auth/register` of `csp-auth-api`) defined in [05-slo.md](05-slo.md): 99.5 % in 30 days,
error budget 0.5 %, which is 75 failed requests per month with the 500 registrations per day I estimated.

## 1. How much budget each alert means before it fires

A burn rate is how many times faster than allowed the budget is being spent: the error ratio divided by the budget (0.005).
The share of the monthly budget spent in a window is `burn rate x window / 720 hours`.

| | Fast burn `AuthRegisterFastBurn` | Slow burn `AuthRegisterSlowBurn` |
|---|---|---|
| Burn rate | 14.4 | 6 |
| Windows | 1 hour, confirmed in the last 5 minutes | 6 hours, confirmed in the last 30 minutes |
| Error ratio that fires it | above 7.2 % (14.4 x 0.005) | above 3 % (6 x 0.005) |
| Budget spent before it fires | 14.4 x 1 / 720 = **2 %** | 6 x 6 / 720 = **5 %** |
| In failed requests (budget 75) | 1.5 requests | 3.75 requests |
| In minutes of full outage (budget 216 min) | 4.3 minutes | 10.8 minutes |
| If nobody acts, the budget is gone in | 720 / 14.4 = 50 hours | 720 / 6 = 120 hours (5 days) |
| Severity | `page`: someone has to look now | `ticket`: it is handled in working hours |
| Owner | team `auth`, Oscar | team `auth`, Oscar |

Why two alerts: the fast one catches a real outage in minutes; the slow one catches a leak that is too small for the fast one
(for example a bug that fails 4 % of the registrations) and that would still finish the budget in five days.
The short window (5 and 30 minutes) is the confirmation: it makes the alert stop when the problem stops, instead of
staying red for the rest of the long window.

## 2. Effect of the low traffic

With about 21 registrations per hour, 7.2 % is one and a half requests, so two failed requests in one hour are enough for
the fast alert. That is what the budget says (2 % of it gone), but it would also fire for two unlucky requests. The rules have a guard,
`increase(auth_register_requests_total[1h]) >= 10` (and 30 requests in 6 hours for the slow one), so a single failure
in a quiet hour does not page. When the real traffic is known, the guard and the 5-minute window should be adjusted.

The alerts for the latency and the freshness SLOs are not written: the PDF asks for the availability SLO, and a third and fourth alert
before knowing the real values would be the start of alert fatigue.

## 3. Two cause-based alerts I will NOT create

| Alert | Why not |
|---|---|
| CPU of the `csp-auth-api` container above 80 % | A registration hashes the password with bcrypt, which is slow and uses CPU on purpose, so a premiere day will show high CPU while every request still answers well. It would page without any user feeling a problem. If the CPU is the real cause, the latency and the 5xx ratio show it first. |
| Any error log line about the database or one restart of the service | `csp-auth-api` is restarted when PostgreSQL is not ready at start (PR #85 of `csp-auth-api`). A single restart or a single failed write is absorbed, so paging on the cause would wake someone for something the system fixes alone. The symptom, 5xx answers to the client, is already covered by the two burn-rate alerts. |

Both are useful as graphs on the dashboard to find the cause during an incident (see the runbook of CP8), but not as pages.
