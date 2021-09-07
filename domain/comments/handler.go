package comments

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"atraf-server/pkg/middleware"
	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
	"atraf-server/pkg/validator"
)

type CreateRequest struct {
	CommentFields         // client-updatable fields
	PostId        uid.UID `json:"post_id" validate:"required"`
	ParentId      uid.UID `json:"parent_id"`
}

type CreateResponse struct {
	CommentId uid.UID `json:"comment_id"`
}

type UpdateResponse struct {
	CommentId uid.UID `json:"comment_id"`
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
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if err := handler.validator.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		commentId, err := handler.service.New(session.AccountId, request.PostId, request.ParentId, request.CommentFields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusCreated, CreateResponse{commentId})
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fields CommentFields

		commentId, err := uid.FromString(chi.URLParam(r, "comment_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if err := handler.validator.Struct(fields); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		commentId, err = handler.service.Update(commentId, fields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusOK, UpdateResponse{commentId})
	}
}

func (handler *Handler) ReadMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		comments, err := handler.service.Comments(postId)
		if err != nil {
			rest.Error(w, http.StatusNotFound)
			return
		}

		rest.Success(w, http.StatusOK, ReadManyResponse{comments})
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
