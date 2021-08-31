package app

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

// DBConnection returns a new sql.DB instance
func DBConnection() (*sql.DB, error) {
	addr := net.JoinHostPort(
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
	)

	config := &mysql.Config{
		Addr:                 addr,
		User:                 os.Getenv("MYSQL_USER"),
		Passwd:               os.Getenv("MYSQL_PASS"),
		DBName:               os.Getenv("MYSQL_NAME"),
		ParseTime:            true,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, err
	}

	// Ping the DB and make sure the server is reachable
	if err := db.Ping(); err != nil {
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

	// Server Info Message
	fmt.Printf("Listening on [%s]...", addr)

	return http.ListenAndServe(addr, handler)
}
