package users

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
	"atraf-server/pkg/validator"
)

type CreateResponse struct {
	UserId uid.UID `json:"user_id"`
}

type ReadOneResponse struct {
	User User `json:"user"`
}

type Handler struct {
	service   *Service
	validator *validator.Validator
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var fields UserFields

		if err := json.NewDecoder(r.Body).Decode(&fields); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if err := handler.validator.Struct(fields); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		userId, err := handler.service.New(fields)
		if err != nil {
			rest.Error(w, http.StatusConflict)
			return
		}

		rest.Success(w, http.StatusCreated, CreateResponse{userId})
	}
}

func (handler *Handler) ReadOne() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := uid.FromString(chi.URLParam(r, "user_id"))
		if err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		user, err := handler.service.User(userId)
		if err != nil {
			rest.Error(w, http.StatusNotFound)
			return
		}

		rest.Success(w, http.StatusOK, ReadOneResponse{user})
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
