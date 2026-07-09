package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

type userRow struct {
	ID       int
	Username string
	Password string
}

func main() {
	password := flag.String("password", "123456", "default password for users with MD5 hashes")
	dsn := flag.String("dsn", "host=127.0.0.1 user=postgres password=td27admin dbname=td27 port=5432 sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	db, err := sql.Open("pgx", *dsn)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping: %v", err)
	}

	rows, err := db.Query("SELECT id, username, password FROM sys_management_user WHERE deleted_at IS NULL")
	if err != nil {
		log.Fatalf("failed to query users: %v", err)
	}
	defer rows.Close()

	var updated int
	for rows.Next() {
		var u userRow
		if err := rows.Scan(&u.ID, &u.Username, &u.Password); err != nil {
			log.Printf("skip row: %v", err)
			continue
		}

		if strings.HasPrefix(u.Password, "$2a$") || strings.HasPrefix(u.Password, "$2b$") {
			fmt.Printf("  [SKIP] %s — already bcrypt\n", u.Username)
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("  [FAIL] %s: %v", u.Username, err)
			continue
		}

		_, err = db.Exec("UPDATE sys_management_user SET password = $1 WHERE id = $2", string(hash), u.ID)
		if err != nil {
			log.Printf("  [FAIL] %s: %v", u.Username, err)
			continue
		}

		fmt.Printf("  [OK]   %s → bcrypt (password: %s)\n", u.Username, *password)
		updated++
	}

	fmt.Printf("\nDone. %d user(s) updated.\n", updated)
}
