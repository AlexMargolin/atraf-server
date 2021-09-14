package users

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
	"atraf-server/pkg/validate"
)

type ReadOneResponse struct {
	User User `json:"user"`
}

type Handler struct {
	service  *Service
	validate *validate.Validate
}

func (handler *Handler) ReadOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := uid.FromString(chi.URLParam(r, "user_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		user, err := handler.service.UserById(userId)
		if err != nil {
			rest.Error(w, http.StatusNotFound)
			return
		}

		rest.Success(w, http.StatusOK, &ReadOneResponse{user})
	}
}

func NewHandler(s *Service, v *validate.Validate) *Handler {
	return &Handler{s, v}
}
