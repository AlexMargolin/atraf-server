package main

import (
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/go-sql-driver/mysql"

	"quotes/config"
)

type Config struct {
	DbHost       string `env:"MYSQL_HOST"`     // Database host name
	DbPort       string `env:"MYSQL_PORT"`     // Database port
	DbUser       string `env:"MYSQL_USER"`     // Database username
	DbPass       string `env:"MYSQL_PASS"`     // Database user password
	DbName       string `env:"MYSQL_NAME"`     // Database DB Name
	ServerPort   string `env:"SERVER_PORT"`    // Server Port
	AuthATSecret string `env:"AUTH_AT_SECRET"` // Authentication Access Token Signing Secret
	AuthRTSecret string `env:"AUTH_RT_SECRET"` // Authentication Refresh Token Signing Secret
}

func main() {
	// Default Configuration
	cfg := &Config{
		DbHost:       "localhost",
		DbPort:       "3306",
		DbUser:       "root",
		DbPass:       "",
		DbName:       "at_db",
		ServerPort:   "8080",
		AuthATSecret: "abcd",
		AuthRTSecret: "dcba",
	}

	if err := config.Marshal(cfg); err != nil {
		log.Fatal(err)
	}

	// MySQL Database connection
	_, err := NewDatabaseConnection(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

// NewDatabaseConnection Returns a new MySql DB Connection or error on failure
func NewDatabaseConnection(c *Config) (*sql.DB, error) {
	// DSN Configuration
	mysqlConfig := &mysql.Config{
		Addr:                 net.JoinHostPort(c.DbHost, c.DbPort),
		User:                 c.DbUser,
		Passwd:               c.DbPass,
		DBName:               c.DbName,
		AllowNativePasswords: true,
	}

	// DSN config string
	dsn := mysqlConfig.FormatDSN()

	db, err := sql.Open("mysql", dsn)
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
