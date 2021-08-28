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
	AccessToken string `json:"accessToken"`
}

type Handler struct {
	service   *Service
	validator *validator.Validator
}

func (handler *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RegisterRequest

		// Decode & Validate JSON
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}

		// Validate request fields
		if err := handler.validator.Struct(request); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity) // 422
			return
		}

		// Register an account
		_, err := handler.service.Register(request.Email, request.Password)
		if err != nil {
			w.WriteHeader(http.StatusConflict) // 409
			return
		}

		// Success
		rest.Response(w, http.StatusCreated) // 201
	}
}

func (handler *Handler) Login(tokenSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request LoginRequest

		// Decode & Validate JSON
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}

		// Validate request
		if err := handler.validator.Struct(request); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity) // 422
			return
		}

		// Account Login
		account, err := handler.service.Login(request.Email, request.Password)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized) // 401
			return
		}

		// Issue Access Token
		accessToken, err := token.New(tokenSecret, account.Id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}

		// Success
		rest.Response(w, http.StatusOK, &LoginResponse{accessToken})
	}
}

// NewHandler returns new account HTTP handler
func NewHandler(service *Service, validator *validator.Validator) *Handler {
	return &Handler{service, validator}
}
