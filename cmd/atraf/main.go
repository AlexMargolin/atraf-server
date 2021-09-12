package main

import (
	"log"

	"github.com/go-chi/chi/v5"

	"atraf-server/services/account"
	"atraf-server/services/comments"
	"atraf-server/services/posts"
	"atraf-server/services/users"

	"atraf-server/app"
	"atraf-server/pkg/middleware"
	"atraf-server/pkg/validate"
)

func main() {
	if err := app.CheckEnvironment(); err != nil {
		log.Fatal(err)
	}

	db, err := app.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	validator := validate.NewValidator()

	accountStorage := account.NewStorage(db)
	accountService := account.NewService(accountStorage)
	accountHandler := account.NewHandler(accountService, validator)

	usersStorage := users.NewStorage(db)
	usersService := users.NewService(usersStorage)
	usersHandler := users.NewHandler(usersService, validator)

	postsStorage := posts.NewStorage(db)
	postsService := posts.NewService(postsStorage)
	postsHandler := posts.NewHandler(postsService, validator)

	commentsStorage := comments.NewStorage(db)
	commentsService := comments.NewService(commentsStorage)
	commentsHandler := comments.NewHandler(commentsService, validator)

	router := chi.NewRouter()
	router.Use(middleware.Cors)
	router.Use(middleware.Options)

	// Unauthenticated Routes (Public)
	// Routes defined under this group do not have access to the Session Context
	router.Group(func(router chi.Router) {
		// Account
		router.Post("/account/register", accountHandler.Register(usersService))
		router.Post("/account/login", accountHandler.Login())
		router.Post("/account/forgot", accountHandler.Forgot())
		router.Patch("/account/reset", accountHandler.Reset())
	})

	// Authenticated Inactive Account Routes (Private)
	// Routes defined under this group have access to the Session Context
	router.Group(func(router chi.Router) {
		router.Use(middleware.Session(false))

		// Account
		router.Post("/account/activate", accountHandler.Activate())
		router.Post("/account/activate/resend", accountHandler.Resend())
	})

	// Authenticated Active Account Routes (Private)
	// Routes defined under this group have access to the Session Context
	router.Group(func(router chi.Router) {
		router.Use(middleware.Session(true))

		// Users
		router.Get("/users/{user_id}", usersHandler.ReadOne())

		// Posts
		router.Post("/posts", postsHandler.Create(usersService))
		router.Put("/posts/{post_id}", postsHandler.Update())
		router.Get("/posts/{post_id}", postsHandler.ReadOne(usersService))
		router.With(middleware.Pagination).Get("/posts", postsHandler.ReadMany(usersService))

		// Comments
		router.Post("/comments", commentsHandler.Create(usersService))
		router.Get("/comments/{source_id}", commentsHandler.ReadMany(usersService))
		router.Put("/comments/{comment_id}", commentsHandler.Update())
	})

	log.Fatal(app.ServeHTTP(router))
}
