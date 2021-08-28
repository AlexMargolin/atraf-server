package app

import (
	"database/sql"
	"net"
	"time"

	"github.com/go-sql-driver/mysql"
)

type SqlConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

func SqlConnection(c *SqlConfig) (*sql.DB, error) {
	config := &mysql.Config{
		Addr:                 net.JoinHostPort(c.Host, c.Port),
		User:                 c.User,
		Passwd:               c.Pass,
		DBName:               c.Name,
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
	db.SetConnMaxLifetime(time.Minute * 3)

	return db, nil
}
