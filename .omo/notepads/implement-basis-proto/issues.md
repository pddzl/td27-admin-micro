
## 2026-05-08: Pre-existing go.sum permission issue

- `go.sum` is owned by `root` (not the current user), causing `go mod tidy` to fail with permission denied.
- This did not block `go build ./...` which passed cleanly.
- Fix would be: `sudo chown $(whoami) go.sum` but sudo requires interactive terminal.
