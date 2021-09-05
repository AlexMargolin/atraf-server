package posts

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"atraf-server/pkg/middleware"
	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
	"atraf-server/pkg/validator"
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
	NextCursor string `json:"next_cursor"`
	Posts      []Post `json:"posts"`
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
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if err := handler.validator.Struct(fields); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		postId, err := handler.service.New(session.AccountId, fields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusCreated, CreateResponse{postId})
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fields PostFields

		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		if err = json.NewDecoder(r.Body).Decode(&fields); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if err = handler.validator.Struct(fields); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		postId, err = handler.service.Update(postId, fields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusOK, UpdateResponse{postId})
	}
}

func (handler *Handler) ReadOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		post, err := handler.service.Post(postId)
		if err != nil {
			rest.Error(w, http.StatusNotFound)
			return
		}

		rest.Success(w, http.StatusOK, ReadOneResponse{post})
	}
}

func (handler *Handler) ReadMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var nextCursor string

		pagination := middleware.GetPaginationContext(r)

		// we want to get the next post as well to determine whether there is a next page.
		limit := pagination.Limit + 1

		posts, err := handler.service.Posts(limit, pagination.Cursor)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if len(posts) == limit {
			nextCursor = posts[len(posts)-2].Id.String()
		}

		rest.Success(w, http.StatusOK, ReadManyResponse{nextCursor, posts})
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
