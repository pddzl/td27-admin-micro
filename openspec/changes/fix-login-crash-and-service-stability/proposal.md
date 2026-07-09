## Why

The service has three stability bugs and the login flow was broken:

1. **Login crashes** — `sql: Scan error on column index 1, name "created_at": unsupported Scan, storing driver.Value type <nil> into type *time.Time`
2. **Basis RPC crashes after ~10 minutes** — SIGSEGV in Casbin enforcer due to nil adapter + auto-load interval
3. **Captcha route mismatch** — route was `POST /captcha` but frontend expects `GET /api/captcha`
4. **Logout 401** — login never succeeds so there's no token to send

## What Changes

- Add COALESCE to userColumns SELECT to handle NULL timestamps from legacy data
- Set `auto-load-interval: 0` in Casbin config to prevent nil adapter SIGSEGV
- Change captcha route from `POST /captcha` to `GET /api/captcha`
- Fix DB row for `admin` user (NULL created_at)

## Impact

No functional changes — fixes crashes and makes login work end-to-end.
