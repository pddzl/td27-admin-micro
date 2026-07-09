## 1. Fix login NULL timestamp scan error

- [x] 1.1 Fix `admin` user row: set `created_at` from NULL to a timestamp in DB
- [x] 1.2 Add `COALESCE(created_at, NOW())` to `userColumns` in user repository

## 2. Fix Casbin SIGSEGV crash

- [x] 2.1 Set `auto-load-interval: 0` in `rpc/basis/etc/basis.yaml` to disable nil adapter auto-reload

## 3. Fix captcha route

- [x] 3.1 Change captcha route from `POST /captcha` to `GET /api/captcha` in handler.go

## 4. Verify

- [ ] 4.1 Restart both services and confirm no login crash
- [ ] 4.2 Confirm basis RPC does not crash after 10+ minutes
- [ ] 4.3 Confirm full login flow: captcha → login with captcha → returns JWT token
- [ ] 4.4 Confirm logout with token returns 200
- [ ] 4.5 Confirm protected endpoint rejects invalidated token after logout
