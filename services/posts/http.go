package posts

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

type CreateRequest = PostFields

type CreateResponse struct {
	PostId uid.UID `json:"post_id"`
}

type UpdateRequest = PostFields

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
	service  *Service
	validate *validate.Validate
}

func (handler *Handler) Create(u *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateRequest
		session := middleware.GetSessionContext(r)

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := handler.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// Dependency(Users)
		user, err := u.UserByAccountId(session.AccountId)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		postId, err := handler.service.NewPost(user.Id, request)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusCreated, &CreateResponse{
			postId,
		})
	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UpdateRequest

		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
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

		if err = handler.service.UpdatePost(postId, request); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
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

		// Dependency(Users)
		user, err := u.UserById(post.UserId)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, &ReadOneResponse{
			post,
			user,
		})
	}
}

func (handler *Handler) ReadMany(u *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pagination := middleware.GetPaginationContext(r)

		posts, err := handler.service.ListPosts(pagination)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		if len(posts) == 0 {
			rest.Error(w, http.StatusNotFound)
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

		userIds := UniqueUserIds(posts)

		// Dependency(Users)
		postsUsers, err := u.UsersByIds(userIds)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, &ReadManyResponse{
			cursor,
			posts,
			postsUsers,
		})
	}
}

func NewHandler(s *Service, v *validate.Validate) *Handler {
	return &Handler{s, v}
}
