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

func main() {
	// Check environment
	if err := app.CheckEnvironment(); err != nil {
		log.Fatal(err)
	}

	// Database(MySQL) Connection
	db, err := app.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	// Struct validator instance
	validate := validator.NewValidator()

	// Account
	accountStorage := account.NewStorage(db)
	accountService := account.NewService(accountStorage)
	accountHandler := account.NewHandler(accountService, validate)

	// Posts
	postsStorage := posts.NewStorage(db)
	postsService := posts.NewService(postsStorage)
	postsHandler := posts.NewHandler(postsService, validate)

	// Comments
	commentsStorage := comments.NewStorage(db)
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
			router.Post("/login", accountHandler.Login())
		})
	})

	// Authenticated Routes (Private)
	// Routes defined under this group have access to the Session Context
	router.Group(func(router chi.Router) {
		router.Use(middleware.Session)

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

	// HTTP Server & Handler
	if err = app.ServeHTTP(router); err != nil {
		log.Fatal(err)
	}
}
