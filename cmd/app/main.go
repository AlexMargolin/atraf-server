package main

import (
	"log"

	"github.com/go-chi/chi/v5"

	"quotes/app"
	"quotes/env"
	"quotes/internal/account"
	"quotes/internal/posts"
	"quotes/pkg/middleware"
	"quotes/pkg/validator"
)

type Config struct {
	DbHost      string `env:"MYSQL_HOST"`   // Database host name
	DbPort      string `env:"MYSQL_PORT"`   // Database port
	DbUser      string `env:"MYSQL_USER"`   // Database username
	DbPass      string `env:"MYSQL_PASS"`   // Database user password
	DbName      string `env:"MYSQL_NAME"`   // Database DB Name
	ServerPort  string `env:"SERVER_PORT"`  // Server Port
	ServerHost  string `env:"SERVER_HOST"`  // Server Port
	TokenSecret string `env:"TOKEN_SECRET"` // Authentication Token Signing Secret
}

func main() {
	config, err := NewConfig(&Config{
		DbHost:      "localhost",
		DbPort:      "3306",
		DbUser:      "root",
		DbPass:      "",
		DbName:      "quotes",
		ServerPort:  "9000",
		TokenSecret: "abcd",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Sql Database
	sqlConfig := &app.SqlConfig{
		Host: config.DbHost,
		Port: config.DbPort,
		User: config.DbUser,
		Pass: config.DbPass,
		Name: config.DbName,
	}
	sql, err := app.SqlConnection(sqlConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Struct validator instance
	validate := validator.NewValidator()

	// Account
	accountStorage := account.NewStorage(sql)
	accountService := account.NewService(accountStorage)
	accountHandler := account.NewHandler(accountService, validate)

	// Posts
	postsStorage := posts.NewStorage(sql)
	postsService := posts.NewService(postsStorage)
	postsHandler := posts.NewHandler(postsService, validate)

	// HTTP Router
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

	// Authenticated Routes (Private)
	// Routes defined under this group have access to the Session Context
	router.Group(func(router chi.Router) {
		router.Use(middleware.Session(config.TokenSecret))

		// Posts
		router.Route("/posts", func(router chi.Router) {
			router.Post("/", postsHandler.Create())                              // Create Post
			router.Put("/{post_id}", postsHandler.Update())                      // Update Post
			router.Get("/{post_id}", postsHandler.ReadOne())                     // Get Post
			router.With(middleware.Pagination).Get("/", postsHandler.ReadMany()) // Get Posts
		})
	})

	// Server
	srvConfig := &app.ServerConfig{
		Host: config.ServerHost,
		Port: config.ServerPort,
	}
	if err = app.RunServer(srvConfig, router); err != nil {
		log.Fatal(err)
	}
}

// NewConfig creates default config and
// assigns defined environment variables
func NewConfig(c *Config) (*Config, error) {
	// Marshal Environment Variables
	err := env.NewDecoder().Marshal(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
