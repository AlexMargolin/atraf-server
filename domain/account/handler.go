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

		// 400
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// 422
		if err := handler.validator.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// 409
		accountId, err := handler.service.Register(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusConflict)
			return
		}

		// 201
		rest.Success(w, http.StatusCreated, RegisterResponse{accountId})
	}
}

func (handler *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request LoginRequest

		// 400
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// 422
		if err := handler.validator.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// 401
		account, err := handler.service.Login(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		// 500
		accessToken, err := token.New(account.Id)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		// 200
		rest.Success(w, http.StatusOK, LoginResponse{accessToken})
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
