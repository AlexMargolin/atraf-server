package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"atraf-server/app"

	"atraf-server/services/account"
	"atraf-server/services/bucket"
	"atraf-server/services/comments"
	"atraf-server/services/posts"
	"atraf-server/services/users"

	"atraf-server/pkg/authentication"
	"atraf-server/pkg/middleware"
	"atraf-server/pkg/validate"
)

func main() {
	if err := app.CheckEnvironment(); err != nil {
		log.Fatal(err)
	}

	sql, err := app.DBConnection()
	if err != nil {
		log.Fatal(err)
	}

	validator := validate.NewValidator()

	bucketStorage := bucket.NewFSBucket()
	bucketService := bucket.NewService(bucketStorage)

	usersStorage := users.NewStorage(sql)
	usersService := users.NewService(usersStorage)
	usersHandler := users.NewHandler(usersService, validator)

	accountStorage := account.NewStorage(sql)
	accountService := account.NewService(accountStorage)
	accountHandler := account.NewHandler(accountService, usersService, validator)

	postsStorage := posts.NewStorage(sql, bucketService)
	postsService := posts.NewService(postsStorage)
	postsHandler := posts.NewHandler(postsService, usersService, validator)

	commentsStorage := comments.NewStorage(sql)
	commentsService := comments.NewService(commentsStorage)
	commentsHandler := comments.NewHandler(commentsService, usersService, validator)

	router := chi.NewRouter()
	router.Use(middleware.Cors)
	router.Use(middleware.Options)

	// health check
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Public Routes
	router.Group(func(router chi.Router) {
		router.Post("/account/register", accountHandler.Register())
		router.Post("/account/login", accountHandler.Login())
		router.Post("/account/forgot", accountHandler.Forgot())
		router.Patch("/account/reset", accountHandler.Reset())
	})

	// Private Routes (unverified users)
	router.Group(func(router chi.Router) {
		router.Use(authentication.Middleware(false))

		router.Patch("/account/activate", accountHandler.Activate())
	})

	// Private Routes (verified users)
	router.Group(func(router chi.Router) {
		router.Use(authentication.Middleware(true))

		// FS Bucket specific file server
		router.Get("/uploads/*", bucketStorage.ServeFiles())

		router.Get("/users/{user_id}", usersHandler.ReadOne())

		router.Post("/posts", postsHandler.Create())
		router.Put("/posts/{post_id}", postsHandler.Update())
		router.Get("/posts/{post_id}", postsHandler.ReadOne())
		router.With(middleware.Pagination).Get("/posts", postsHandler.ReadMany())

		router.Post("/comments", commentsHandler.Create())
		router.Get("/comments/{source_id}", commentsHandler.ReadMany())
		router.Put("/comments/{comment_id}", commentsHandler.Update())
	})

	if err = app.ServeHTTP(router); err != nil {
		log.Fatal()
	}
}
