package account

import (
	"encoding/json"
	"net/http"

	"quotes/pkg/rest"
	"quotes/pkg/token"
	"quotes/pkg/validator"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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
			rest.Error(w, http.StatusBadRequest) // 400
			return
		}

		if err := handler.validator.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity) // 422
			return
		}

		_, err := handler.service.Register(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusConflict) // 409
			return
		}

		rest.Success(w, http.StatusCreated) // 201
	}
}

func (handler *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusBadRequest) // 400
			return
		}

		if err := handler.validator.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity) // 422
			return
		}

		account, err := handler.service.Login(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized) // 401
			return
		}

		accessToken, err := token.New(account.Id)
		if err != nil {
			rest.Error(w, http.StatusInternalServerError) // 500
			return
		}

		rest.Success(w, http.StatusOK, LoginResponse{accessToken})
	}
}

func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
