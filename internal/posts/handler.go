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
	Id uid.UID `json:"id"`
}

type UpdateResponse struct {
	Id uid.UID `json:"id"`
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

		// Session Context
		session := middleware.GetSessionContext(r)

		// Decode & Validate JSON
		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			http.Error(w, "", http.StatusBadRequest) // 400
			return
		}

		// Validate Payload
		if err := handler.validator.Struct(fields); err != nil {
			http.Error(w, "", http.StatusUnprocessableEntity) // 422
			return
		}

		// Create a new Storage Entry
		postId, err := handler.service.New(session.AccountId, fields)
		if err != nil {

			http.Error(w, "", http.StatusBadRequest) // 400
			return
		}

		// Success
		rest.Response(w, http.StatusCreated, &CreateResponse{postId}) // 201
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fields PostFields

		// uid from string
		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "", http.StatusUnprocessableEntity) // 422
			return
		}

		// Decode & Validate JSON
		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			http.Error(w, "", http.StatusBadRequest) // 400
			return
		}

		// Validate Payload
		if err := handler.validator.Struct(fields); err != nil {
			http.Error(w, "", http.StatusUnprocessableEntity) // 422
			return
		}

		// Update existing Storage Entry
		postId, err = handler.service.Update(postId, fields)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest) // 400
			return
		}

		// API Success
		rest.Response(w, http.StatusOK, &UpdateResponse{postId}) // 200
	}
}

func (handler *Handler) ReadOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			http.Error(w, "", http.StatusUnprocessableEntity) // 422
			return
		}

		post, err := handler.service.Post(postId)
		if err != nil {
			http.Error(w, "", http.StatusNotFound) // 404
			return
		}

		// Success
		rest.Response(w, http.StatusOK, &ReadOneResponse{post}) // 200
	}
}

func (handler *Handler) ReadMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Pagination Context
		pagination := middleware.GetPaginationContext(r)

		// Pagination Params validation
		if pagination.PageNum < 1 || pagination.PerPage < 1 {
			http.Error(w, "", http.StatusUnprocessableEntity) // 422
			return
		}

		// Total posts count
		total, err := handler.service.Total()
		if err != nil {
			http.Error(w, "", http.StatusBadRequest) // 400
			return
		}

		// Determines whether current page will return any posts.
		// Otherwise, we return an empty posts array.
		if !middleware.HasCurrentPage(pagination, total) {
			rest.Response(w, http.StatusOK, &ReadManyResponse{Total: total}) // 200
			return
		}

		// Fetch Posts list
		posts, err := handler.service.Posts(pagination.PageNum, pagination.PerPage)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest) // 400
			return
		}

		// Determine whether next page will return any posts
		hasNext := middleware.HasNextPage(pagination, total)

		// API Success
		rest.Response(w, http.StatusOK, &ReadManyResponse{total, hasNext, posts}) // 200
	}
}

// NewHandler returns new Posts HTTP handler
func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
