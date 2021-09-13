package comments

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"atraf-server/pkg/middleware"
	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
	"atraf-server/pkg/validate"
	"atraf-server/services/users"
)

type CreateRequest struct {
	CommentFields
	SourceId uid.UID `json:"source_id" validate:"required"`
	ParentId uid.UID `json:"parent_id"`
}

type CreateResponse struct {
	CommentId uid.UID `json:"comment_id"`
}

type UpdateRequest = CommentFields

type ReadManyResponse struct {
	Comments []Comment    `json:"comments"`
	Users    []users.User `json:"users"`
}

type Handler struct {
	service  *Service
	validate *validate.Validate
}

func (handler *Handler) Create(u *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateRequest
		auth := middleware.GetAuthContext(r)

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := handler.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// Dependency(Users)
		user, err := u.UserByAccountId(auth.AccountId)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		commentId, err := handler.service.NewComment(user.Id, request.SourceId, request.ParentId, request.CommentFields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusCreated, &CreateResponse{
			commentId,
		})
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UpdateRequest

		commentId, err := uid.FromString(chi.URLParam(r, "comment_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err = handler.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		if err = handler.service.UpdateComment(commentId, request); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func (handler *Handler) ReadMany(u *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sourceId, err := uid.FromString(chi.URLParam(r, "source_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		comments, err := handler.service.CommentsBySourceId(sourceId)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		if len(comments) == 0 {
			rest.Success(w, http.StatusOK, &ReadManyResponse{
				[]Comment{},
				[]users.User{},
			})
			return
		}

		userIds := UniqueUserIds(comments)

		// Dependency(Users)
		__users, err := u.UsersByIds(userIds)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, &ReadManyResponse{
			comments,
			__users,
		})
	}
}

func NewHandler(s *Service, v *validate.Validate) *Handler {
	return &Handler{s, v}
}
