package main

import (
	"log"

	"github.com/go-chi/chi/v5"

	"quotes/domain/account"
	"quotes/domain/comments"
	"quotes/domain/posts"

	"quotes/app"
	"quotes/pkg/middleware"
	"quotes/pkg/validator"
)

func main() {
	if err := app.CheckEnvironment(); err != nil {
		log.Fatal(err)
	}

	db, err := app.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.NewValidator()

	accountStorage := account.NewStorage(db)
	accountService := account.NewService(accountStorage)
	accountHandler := account.NewHandler(accountService, validate)

	postsStorage := posts.NewStorage(db)
	postsService := posts.NewService(postsStorage)
	postsHandler := posts.NewHandler(postsService, validate)

	commentsStorage := comments.NewStorage(db)
	commentsService := comments.NewService(commentsStorage)
	commentsHandler := comments.NewHandler(commentsService, validate)

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

		router.Route("/posts", func(router chi.Router) {
			router.Post("/", postsHandler.Create())
			router.Put("/{post_id}", postsHandler.Update())
			router.Get("/{post_id}", postsHandler.ReadOne())
			router.With(middleware.Pagination).Get("/", postsHandler.ReadMany())
		})

		router.Route("/comments", func(router chi.Router) {
			router.Post("/", commentsHandler.Create())
			router.Get("/{post_id}", commentsHandler.ReadMany())
			router.Put("/{comment_id}", commentsHandler.Update())
		})
	})

	log.Fatal(app.ServeHTTP(router))
}
