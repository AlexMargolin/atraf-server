package comments

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"atraf-server/pkg/middleware"
	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
	"atraf-server/pkg/validator"
	"atraf-server/services/users"
)

type CreateRequest struct {
	CommentFields         // client-updatable fields
	SourceId      uid.UID `json:"source_id" validate:"required"`
	ParentId      uid.UID `json:"parent_id"`
}

type CreateResponse struct {
	CommentId uid.UID `json:"comment_id"`
}

type UpdateResponse struct {
	CommentId uid.UID `json:"comment_id"`
}

type ReadManyResponse struct {
	Comments []Comment    `json:"comments"`
	Users    []users.User `json:"users"`
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

		commentId, err := handler.service.NewComment(session.UserId, request.SourceId, request.ParentId, request.CommentFields)
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

		commentId, err = handler.service.UpdateComment(commentId, fields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusOK, UpdateResponse{commentId})
	}
}

// ReadMany Depends on: Users
func (handler *Handler) ReadMany(u *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sourceId, err := uid.FromString(chi.URLParam(r, "source_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		comments, err := handler.service.CommentsBySourceId(sourceId)
		if err != nil {
			rest.Error(w, http.StatusNotFound)
			return
		}

		if len(comments) == 0 {
			rest.Success(w, http.StatusOK, ReadManyResponse{})
			return
		}

		commentsUserIds := UniqueUserIds(comments)

		// TODO replace with endpoint
		__users, err := u.UsersByIds(commentsUserIds)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, ReadManyResponse{comments, __users})
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
