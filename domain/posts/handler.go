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
	Posts []Post `json:"posts"`
}

type Handler struct {
	service   *Service
	validator *validator.Validator
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fields PostFields

		session := middleware.GetSessionContext(r)

		// 400
		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// 422
		if err := handler.validator.Struct(fields); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// 400
		postId, err := handler.service.New(session.AccountId, fields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// 201
		rest.Success(w, http.StatusCreated, CreateResponse{postId})
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fields PostFields

		// 422
		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// 400
		if err = json.NewDecoder(r.Body).Decode(&fields); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// 422
		if err = handler.validator.Struct(fields); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// 400
		postId, err = handler.service.Update(postId, fields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// 200
		rest.Success(w, http.StatusOK, UpdateResponse{postId})
	}
}

func (handler *Handler) ReadOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 422
		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// 404
		post, err := handler.service.Post(postId)
		if err != nil {
			rest.Error(w, http.StatusNotFound)
			return
		}

		// 200
		rest.Success(w, http.StatusOK, ReadOneResponse{post})
	}
}

func (handler *Handler) ReadMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := middleware.GetPaginationContext(r)

		// 400
		posts, err := handler.service.Posts(pagination.Limit, pagination.Cursor)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// 200
		rest.Success(w, http.StatusOK, ReadManyResponse{posts})
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
