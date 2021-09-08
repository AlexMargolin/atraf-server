package account

import (
	"encoding/json"
	"net/http"

	"atraf-server/services/users"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/uid"
	"atraf-server/pkg/validator"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	UserId uid.UID `json:"user_id"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type ForgotRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type Handler struct {
	service   *Service
	validator *validator.Validator
}

// Register Depends on: Users
func (handler *Handler) Register(u *users.Service) http.HandlerFunc {
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

		// DOMAIN Dependency (Users)
		__userId, err := u.NewUser(accountId, users.UserFields{Email: request.Email})
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusCreated, RegisterResponse{__userId})
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

		token, err := handler.service.Login(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		rest.Success(w, http.StatusOK, LoginResponse{token})
	}
}

// Forgot attempts to locate an account using the request email.
// once located, a password reset email is sent containing a signed JWT with the account id.
// returns the same status code regardless of whether the account exists or not.
func (handler *Handler) Forgot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ForgotRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if err := handler.validator.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		if err := handler.service.Forgot(request.Email); err != nil {
			rest.Success(w, http.StatusNoContent)
			return
		}

		rest.Success(w, http.StatusNoContent)
	}
}

func NewHandler(service *Service, v *validator.Validator) *Handler {
	return &Handler{service, v}
}
