package main

import (
	"log"

	"github.com/go-chi/chi/v5"

	"quotes/app"
	"quotes/internal/account"
	"quotes/internal/comments"
	"quotes/internal/posts"
	"quotes/pkg/middleware"
	"quotes/pkg/validator"
)

type Config struct {
	DbHost      string `env:"DB_HOST"`      // Database host name
	DbPort      string `env:"DB_PORT" `     // Database port
	DbUser      string `env:"DB_USER"`      // Database username
	DbPass      string `env:"DB_PASS"`      // Database user password
	DbName      string `env:"DB_NAME"`      // Database DB Name
	ServerPort  string `env:"SERVER_PORT"`  // Server Port
	ServerHost  string `env:"SERVER_HOST"`  // Server Port
	TokenSecret string `env:"TOKEN_SECRET"` // Authentication Token Signing Secret
}

func main() {
	// App Config
	var config Config
	if err := app.NewConfig().Decode(&config); err != nil {
		log.Fatal(err)
	}

	// Database
	sql, err := app.SqlConnection(&app.SqlConfig{
		Host: config.DbHost,
		Port: config.DbPort,
		User: config.DbUser,
		Pass: config.DbPass,
		Name: config.DbName,
	})
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

	// Comments
	commentsStorage := comments.NewStorage(sql)
	commentsService := comments.NewService(commentsStorage)
	commentsHandler := comments.NewHandler(commentsService, validate)

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
			router.Post("/", postsHandler.Create())
			router.Put("/{post_id}", postsHandler.Update())
			router.Get("/{post_id}", postsHandler.ReadOne())
			router.With(middleware.Pagination).Get("/", postsHandler.ReadMany())
		})

		// Comments
		router.Route("/comments", func(router chi.Router) {
			router.Post("/", commentsHandler.Create())
			router.Get("/{post_id}", commentsHandler.ReadMany())
			router.Put("/{comment_id}", commentsHandler.Update())
		})
	})

	// Server
	err = app.Run(&app.ServerConfig{
		Host:    config.ServerHost,
		Port:    config.ServerPort,
		Handler: router,
	})
	if err != nil {
		log.Fatal(err)
	}
}
