package account

import (
	"encoding/json"
	"net/http"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/token"
	"atraf-server/pkg/uid"
	"atraf-server/pkg/validator"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	AccountId uid.UID `json:"account_id"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type Handler struct {
	service   *Service
	validator *validator.Validator
}

func (handler *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if err := handler.validator.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		accountId, err := handler.service.Register(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusConflict)
			return
		}

		rest.Success(w, http.StatusCreated, RegisterResponse{accountId})
	}
}

func (handler *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if err := handler.validator.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		account, err := handler.service.Login(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		accessToken, err := token.New(account.Id)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, LoginResponse{accessToken})
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
