package comments

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"quotes/pkg/middleware"
	"quotes/pkg/rest"
	"quotes/pkg/uid"
	"quotes/pkg/validator"
)

type CreateRequest struct {
	CommentFields         // client-updatable fields
	PostId        uid.UID `json:"post_id" validate:"required"`
	ParentId      uid.UID `json:"parent_id"`
}

type CreateResponse struct {
	CommentId uid.UID `json:"id"`
}

type UpdateResponse struct {
	CommentId uid.UID `json:"id"`
}

type ReadManyResponse struct {
	Comments []Comment `json:"comments"`
}

type Handler struct {
	service   *Service
	validator *validator.Validator
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateRequest

		session := middleware.GetSessionContext(r)

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}

		if err := handler.validator.Struct(request); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity) // 422
			return
		}

		commentId, err := handler.service.New(session.AccountId, request.PostId, request.ParentId, request.CommentFields)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}

		rest.Response(w, http.StatusCreated, &CreateResponse{commentId}) // 201
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fields CommentFields

		commentId, err := uid.FromString(chi.URLParam(r, "comment_id"))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity) // 422
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}

		if err := handler.validator.Struct(fields); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity) // 422
			return
		}

		commentId, err = handler.service.Update(commentId, fields)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}

		rest.Response(w, http.StatusOK, &UpdateResponse{commentId}) // 200
	}
}

func (handler *Handler) ReadMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity) // 422
			return
		}

		comments, err := handler.service.Comments(postId)
		if err != nil {
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}

		rest.Response(w, http.StatusOK, &ReadManyResponse{comments}) // 200
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
