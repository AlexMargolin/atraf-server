package comments

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"atraf-server/pkg/authentication"
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
	Comment Comment `json:"comment"`
}

type UpdateRequest = CommentFields

type ReadManyResponse struct {
	Comments []Comment    `json:"comments"`
	Users    []users.User `json:"users"`
}

type Handler struct {
	service  *Service
	users    *users.Service
	validate *validate.Validate
}

func (h Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateRequest
		auth := authentication.Context(r)

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := h.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// Dependency(Users)
		user, err := h.users.UserByAccountId(auth.AccountId)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		comment, err := h.service.NewComment(user.Id, request.SourceId, request.ParentId, request.CommentFields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusCreated, &CreateResponse{
			comment,
		})
	}
}

func (h Handler) Update() http.HandlerFunc {
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

		if err = h.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		if err = h.service.UpdateComment(commentId, request); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func (h Handler) ReadMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sourceId, err := uid.FromString(chi.URLParam(r, "source_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		comments, err := h.service.CommentsBySourceId(sourceId)
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
		__users, err := h.users.UsersByIds(userIds)
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

func NewHandler(s *Service, u *users.Service, v *validate.Validate) *Handler {
	return &Handler{s, u, v}
}
