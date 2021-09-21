package posts

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"atraf-server/services/users"

	"atraf-server/pkg/authentication"
	"atraf-server/pkg/middleware"
	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
	"atraf-server/pkg/validate"
)

const (
	AttachmentMaxSize = 10 * 1024 * 1024 // 10MB
	AttachmentFormKey = "attachment"
)

type CreateRequest = Fields

type CreateResponse struct {
	PostId uid.UID `json:"post_id"`
}

type UpdateRequest = Fields

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
	users    *users.Service
	validate *validate.Validate
}

func (h Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := authentication.Context(r)

		// set max request size
		r.Body = http.MaxBytesReader(w, r.Body, AttachmentMaxSize)

		// set max size allowed before writing to the filesystem.
		if err := r.ParseMultipartForm(AttachmentMaxSize); err != nil {
			rest.Error(w, err, http.StatusRequestEntityTooLarge)
			return
		}
		defer r.Body.Close()

		file, _, err := r.FormFile(AttachmentFormKey)
		if err != nil {
			rest.Error(w, err, http.StatusUnprocessableEntity)
			return
		}
		defer file.Close()

		request := &CreateRequest{
			Title: r.FormValue("title"),
			Body:  r.FormValue("body"),
			File:  file,
		}

		if err = h.validate.Struct(request); err != nil {
			rest.Error(w, err, http.StatusUnprocessableEntity)
			return
		}

		// Dependency(Users)
		__user, err := h.users.UserByAccountId(auth.AccountId)
		if err != nil {
			rest.Error(w, err, http.StatusInternalServerError)
			return
		}

		postId, err := h.service.NewPost(__user.Id, request)
		if err != nil {
			rest.Error(w, err, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusCreated, &CreateResponse{
			postId,
		})
	}
}

func (h Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request UpdateRequest

		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, err, http.StatusUnprocessableEntity)
			return
		}

		if err = json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, err, http.StatusUnsupportedMediaType)
			return
		}

		if err = h.validate.Struct(request); err != nil {
			rest.Error(w, err, http.StatusUnprocessableEntity)
			return
		}

		if err = h.service.UpdatePost(postId, &request); err != nil {
			rest.Error(w, err, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func (h Handler) ReadOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postId, err := uid.FromString(chi.URLParam(r, "post_id"))
		if err != nil {
			rest.Error(w, err, http.StatusUnprocessableEntity)
			return
		}

		post, err := h.service.PostById(postId)
		if err != nil {
			rest.Error(w, err, http.StatusNotFound)
			return
		}

		// Dependency(Users)
		__user, err := h.users.UserById(post.UserId)
		if err != nil {
			rest.Error(w, err, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, &ReadOneResponse{
			post,
			__user,
		})
	}
}

func (h Handler) ReadMany() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cursor string

		pagination := middleware.GetPaginationContext(r)

		// we add additional post in order to determine if there is another
		// page available for pagination
		pagination.Limit++

		posts, err := h.service.Posts(pagination)
		if err != nil {
			rest.Error(w, err, http.StatusInternalServerError)
			return
		}

		if len(posts) == 0 {
			rest.Error(w, err, http.StatusNotFound)
			return
		}

		// if both are equal, it means there are more posts
		// than originally queried by the client.
		// in which case, a pagination cursor is added to the response.
		if len(posts) == pagination.Limit {
			// remove the additional post from the posts result
			posts = posts[:len(posts)-1]
			lastPost := posts[len(posts)-1]

			cursor, err = middleware.EncodeCursor(&middleware.Cursor{
				Key:   lastPost.Id,
				Value: lastPost.CreatedAt,
			})

			if err != nil {
				rest.Error(w, err, http.StatusInternalServerError)
				return
			}
		}

		userIds := UniqueUserIds(posts)

		// Dependency(Users)
		__users, err := h.users.UsersByIds(userIds)
		if err != nil {
			rest.Error(w, err, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, &ReadManyResponse{
			cursor,
			posts,
			__users,
		})
	}
}

func NewHandler(s *Service, u *users.Service, v *validate.Validate) *Handler {
	return &Handler{s, u, v}
}
