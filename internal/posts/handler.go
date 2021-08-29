package posts

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"quotes/pkg/middleware"
	"quotes/pkg/rest"
	"quotes/pkg/uid"
	"quotes/pkg/validator"
)

type CreateResponse struct {
	PostId uid.UID `json:"id"`
}

type UpdateResponse struct {
	PostId uid.UID `json:"id"`
}

type ReadOneResponse struct {
	Post Post `json:"post"`
}

type ReadManyResponse struct {
	Total   int    `json:"total"`
	HasNext bool   `json:"has_next"`
	Posts   []Post `json:"posts"`
}

type Handler struct {
	service   *Service
	validator *validator.Validator
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fields PostFields

		session := middleware.GetSessionContext(r)

		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			rest.Error(w, http.StatusBadRequest) // 400
			return
		}

		if err := handler.validator.Struct(fields); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity) // 422
			return
		}

		postId, err := handler.service.New(session.AccountId, fields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest) // 400
			return
		}

		rest.Success(w, http.StatusCreated, CreateResponse{postId}) // 201
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fields PostFields

		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity) // 422
			return
		}

		if err = json.NewDecoder(r.Body).Decode(&fields); err != nil {
			rest.Error(w, http.StatusBadRequest) // 400
			return
		}

		if err = handler.validator.Struct(fields); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity) // 422
			return
		}

		postId, err = handler.service.Update(postId, fields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest) // 400
			return
		}

		rest.Success(w, http.StatusOK, UpdateResponse{postId}) // 200
	}
}

func (handler *Handler) ReadOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity) // 422
			return
		}

		post, err := handler.service.Post(postId)
		if err != nil {
			rest.Error(w, http.StatusNotFound) // 404
			return
		}

		rest.Success(w, http.StatusOK, ReadOneResponse{post}) // 200
	}
}

func (handler *Handler) ReadMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := middleware.GetPaginationContext(r)

		if pagination.PageNum < 1 || pagination.PerPage < 1 {
			rest.Error(w, http.StatusUnprocessableEntity) // 422
			return
		}

		total, err := handler.service.Total()
		if err != nil {
			rest.Error(w, http.StatusBadRequest) // 400
			return
		}

		// Determines whether current page will return any posts.
		// Otherwise, we return an empty posts array.
		if !middleware.HasCurrentPage(pagination, total) {
			rest.Success(w, http.StatusOK, ReadManyResponse{Total: total}) // 200
			return
		}

		posts, err := handler.service.Posts(pagination.PageNum, pagination.PerPage)
		if err != nil {
			rest.Error(w, http.StatusBadRequest) // 400
			return
		}

		// Determine whether next page will return any posts
		hasNext := middleware.HasNextPage(pagination, total)

		rest.Success(w, http.StatusOK, ReadManyResponse{total, hasNext, posts}) // 200
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
