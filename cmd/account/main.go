package main

import (
	"database/sql"
	"log"
	"net"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"

	"quotes/env"
	"quotes/internal/account"
	"quotes/pkg/validator"
)

type Config struct {
	DbHost      string `env:"MYSQL_HOST"`   // Database host name
	DbPort      string `env:"MYSQL_PORT"`   // Database port
	DbUser      string `env:"MYSQL_USER"`   // Database username
	DbPass      string `env:"MYSQL_PASS"`   // Database user password
	DbName      string `env:"MYSQL_NAME"`   // Database DB Name
	ServerPort  string `env:"SERVER_PORT"`  // Server Port
	TokenSecret string `env:"TOKEN_SECRET"` // Authentication Token Signing Secret
}

func main() {
	// Global configuration
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Database configuration
	db, err := NewDbConnection(config)
	if err != nil {
		log.Fatal(err)
	}

	// HTTP Server
	server := NewServer(config)

	// Struct validator instance
	validate := validator.NewValidator()

	// Account
	accountStorage := account.NewStorage(db)
	accountService := account.NewService(accountStorage)
	accountHandler := account.NewHandler(accountService, validate)

	// HTTP router
	router := chi.NewRouter()

	// Unauthenticated Routes (Public)
	// Routes defined under this group do not have access to the Session Context
	router.Group(func(router chi.Router) {
		// Account
		router.Route("/account", func(router chi.Router) {
			router.Post("/register", accountHandler.Register())
			router.Post("/login", accountHandler.Login(config.TokenSecret))
		})
	})

	// Run http server
	if err := server.Run(router); err != nil {
		log.Fatal(err)
	}
}

// NewConfig Creates new Config and sets default values
func NewConfig() (*Config, error) {
	config := &Config{
		DbHost:      "localhost",
		DbPort:      "3306",
		DbUser:      "root",
		DbPass:      "",
		DbName:      "account_db",
		ServerPort:  "8080",
		TokenSecret: "abcd",
	}

	// Parse Environment Variables
	if err := env.NewDecoder().Marshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// NewDbConnection Returns a new MySql DB Connection or error on failure
func NewDbConnection(c *Config) (*sql.DB, error) {
	config := &mysql.Config{
		Addr:                 net.JoinHostPort(c.DbHost, c.DbPort),
		User:                 c.DbUser,
		Passwd:               c.DbPass,
		DBName:               c.DbName,
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
