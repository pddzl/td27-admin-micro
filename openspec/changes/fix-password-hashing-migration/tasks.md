## 1. Create migration script

- [x] 1.1 Create `scripts/migrate_passwords.go` — reads all users with MD5 hashes, replaces with bcrypt

## 2. Run migration

- [x] 2.1 Run `go run scripts/migrate_passwords.go -password 123456`
- [x] 2.2 Verify: all 3 users now have bcrypt hashes
