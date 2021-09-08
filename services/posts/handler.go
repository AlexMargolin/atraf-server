package posts

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

type CreateResponse struct {
	PostId uid.UID `json:"post_id"`
}

type UpdateResponse struct {
	PostId uid.UID `json:"post_id"`
}

type ReadOneResponse struct {
	Post Post       `json:"post"`
	User users.User `json:"user"`
}

type ReadManyResponse struct {
	Cursor string       `json:"cursor"`
	Posts  []Post       `json:"posts"`
	Users  []users.User `json:"users"`
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

		postId, err := handler.service.NewPost(session.UserId, fields)
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

		postId, err = handler.service.UpdatePost(postId, fields)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusOK, UpdateResponse{postId})
	}
}

func (handler *Handler) ReadOne(u *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		post, err := handler.service.PostById(postId)
		if err != nil {
			rest.Error(w, http.StatusNotFound)
			return
		}

		// DOMAIN Dependency (Users)
		__user, err := u.UserById(post.UserId)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, ReadOneResponse{post, __user})
	}
}

func (handler *Handler) ReadMany(u *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := middleware.GetPaginationContext(r)

		posts, err := handler.service.ListPosts(pagination)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if len(posts) == 0 {
			rest.Success(w, http.StatusOK, ReadManyResponse{})
			return
		}

		lastPost := posts[len(posts)-1]
		cursor, err := middleware.EncodeCursor(&middleware.Cursor{
			Key:   lastPost.Id,
			Value: lastPost.CreatedAt,
		})
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		postsUserIds := UniqueUserIds(posts)

		// DOMAIN Dependency (Users)
		__users, err := u.UsersByIds(postsUserIds)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, ReadManyResponse{cursor, posts, __users})
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
