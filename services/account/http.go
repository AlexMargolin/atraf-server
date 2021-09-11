package account

import (
	"encoding/json"
	"net/http"
	"os"

	"atraf-server/pkg/token"
	"atraf-server/pkg/uid"
	"atraf-server/services/users"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/validate"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ActivateRequest struct {
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

type ResetRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type Handler struct {
	service  *Service
	validate *validate.Validate
}

func (handler *Handler) Register(u *users.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := handler.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		accountId, err := handler.service.Register(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusConflict)
			return
		}

		// Dependency(Users) : could be a webhook.
		if err = u.NewUser(accountId, users.UserFields{Email: request.Email}); err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusCreated, nil)
	}
}

func (handler *Handler) Activate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO implement...
	}
}

func (handler *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := handler.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		accessToken, err := handler.service.Login(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		rest.Success(w, http.StatusOK, &LoginResponse{
			accessToken,
		})
	}
}

func (handler *Handler) Forgot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ForgotRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := handler.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// Whether an account can be found or not, a "successful" response is returned.
		if err := handler.service.Forgot(request.Email); err != nil {
			rest.Success(w, http.StatusNoContent, nil)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func (handler *Handler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ResetRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := handler.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		claims, err := token.Verify(os.Getenv("RESET_TOKEN_SECRET"), request.Token)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		accountId, err := uid.FromString(claims.Subject)
		if err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		if err = handler.service.UpdatePassword(accountId, request.NewPassword); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func NewHandler(s *Service, v *validate.Validate) *Handler {
	return &Handler{s, v}
}
