package app

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func DBConnection() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require timezone=utc",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Ping the DB and make sure the server is reachable
	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Second * 10)

	return db, nil
}

// ServeHTTP serves an unencrypted HTTP server
func ServeHTTP(handler http.Handler) error {
	addr := net.JoinHostPort(
		os.Getenv("SERVER_HOST"),
		os.Getenv("SERVER_PORT"),
	)

	fmt.Printf("Listening on [%s]...", addr)

	return http.ListenAndServe(addr, handler)
}
