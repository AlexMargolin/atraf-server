package main

import (
	"log"

	"github.com/go-chi/chi/v5"

	"atraf-server/domain/account"
	"atraf-server/domain/comments"
	"atraf-server/domain/posts"
	"atraf-server/domain/users"

	"atraf-server/app"
	"atraf-server/pkg/middleware"
	"atraf-server/pkg/validator"
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

	usersStorage := users.NewStorage(db)
	usersService := users.NewService(usersStorage)
	usersHandler := users.NewHandler(usersService, validate)

	postsStorage := posts.NewStorage(db)
	postsService := posts.NewService(postsStorage)
	postsHandler := posts.NewHandler(postsService, validate)

	commentsStorage := comments.NewStorage(db)
	commentsService := comments.NewService(commentsStorage)
	commentsHandler := comments.NewHandler(commentsService, validate)

	router := chi.NewRouter()
	router.Use(middleware.Cors)
	router.Use(middleware.Options)

	// Unauthenticated Routes (Public)
	// Routes defined under this group do not have access to the Session Context
	router.Group(func(router chi.Router) {
		// Account
		router.Post("/account/register", accountHandler.Register(usersService))
		router.Post("/account/login", accountHandler.Login())
	})

	// Authenticated Routes (Private)
	// Routes defined under this group have access to the Session Context
	router.Group(func(router chi.Router) {
		router.Use(middleware.Session(usersService))

		// Users
		router.Get("/users/{user_id}", usersHandler.ReadOne())

		// Posts
		router.Post("/posts", postsHandler.Create())
		router.Put("/posts/{post_id}", postsHandler.Update())
		router.Get("/posts/{post_id}", postsHandler.ReadOne())
		router.With(middleware.Pagination).Get("/posts", postsHandler.ReadMany())

		// Comments
		router.Post("/comments", commentsHandler.Create())
		router.Get("/comments/{source_id}", commentsHandler.ReadMany(usersService))
		router.Put("/comments/{comment_id}", commentsHandler.Update())
	})

	log.Fatal(app.ServeHTTP(router))
}
