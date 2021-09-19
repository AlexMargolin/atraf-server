package account

import (
	"encoding/json"
	"net/http"

	"atraf-server/services/users"

	"atraf-server/pkg/authentication"
	"atraf-server/pkg/rest"
	"atraf-server/pkg/token"
	"atraf-server/pkg/validate"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ActivateRequest struct {
	Code string `json:"code" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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
	users    *users.Service
	validate *validate.Validate
}

func (h Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request RegisterRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := h.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		account, err := h.service.Register(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusConflict)
			return
		}

		// Dependency(Users)
		// This could be a webhook
		if h.users.NewUser(account.Id, users.UserFields{Email: account.Email}) != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		if err = authentication.SetCookie(w, account.Id, false); err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func (h Handler) Activate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ActivateRequest
		auth := authentication.Context(r)

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := h.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		if err := h.service.Activate(auth.AccountId, request.Code); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// Issue New Access Token
		if err := authentication.SetCookie(w, auth.AccountId, true); err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func (h Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := h.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		account, err := h.service.Login(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		if err = authentication.SetCookie(w, account.Id, account.Active); err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func (h Handler) Forgot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ForgotRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := h.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		// Whether an account can be found or not, a "successful" response is returned.
		if err := h.service.Forgot(request.Email); err != nil {
			rest.Success(w, http.StatusNoContent, nil)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func (h Handler) Reset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ResetRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := h.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		t, err := token.VerifyResetToken(request.Token)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		if err = h.service.UpdatePassword(t.AccountId, request.NewPassword); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func NewHandler(s *Service, u *users.Service, v *validate.Validate) *Handler {
	return &Handler{s, u, v}
}
