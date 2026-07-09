## Why

The original `td27-admin/server` used **MD5** for password hashing. This rewrite uses **bcrypt** (`golang.org/x/crypto/bcrypt`). All existing user passwords in the DB are MD5 hashes, which bcrypt cannot verify. Every login attempt fails with "invalid username or password" regardless of the correct password.

The admin user was hot-fixed, but other users (edit, x1) are still broken.

## What Changes

- **Migration script**: One-time Go script that reads all users, re-hashes their passwords with bcrypt
- **Or** backward-compatible `VerifyPassword`: try bcrypt first, fall back to MD5 if hash is MD5-format, and upgrade to bcrypt on successful login
- Run migration against all existing users

## Impact

- All existing users can log in again
- New users created via `/api/user/create` already use bcrypt (correct)
- No schema changes
