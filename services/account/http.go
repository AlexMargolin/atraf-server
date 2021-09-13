package account

import (
	"encoding/json"
	"net/http"

	"atraf-server/pkg/middleware"
	"atraf-server/pkg/session"
	"atraf-server/pkg/token"
	"atraf-server/services/users"

	"atraf-server/pkg/rest"
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

		account, err := handler.service.Register(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusConflict)
			return
		}

		// Dependency(Users)
		// TODO should be part of a transaction?
		if u.NewUser(account.Id, users.UserFields{Email: account.Email}) != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		if err = session.SetCookie(w, account.Id, account.Active); err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func (handler *Handler) Activate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ActivateRequest
		auth := middleware.GetAuthContext(r)

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := handler.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		if err := handler.service.Activate(auth.AccountId, request.Code); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// Issue New Access Token
		if err := session.SetCookie(w, auth.AccountId, true); err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
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

		account, err := handler.service.Login(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		if err = session.SetCookie(w, account.Id, account.Active); err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
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

		t, err := token.VerifyResetToken(request.Token)
		if err != nil {
			rest.Error(w, http.StatusUnauthorized)
			return
		}

		if err = handler.service.UpdatePassword(t.AccountId, request.NewPassword); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		rest.Success(w, http.StatusNoContent, nil)
	}
}

func NewHandler(s *Service, v *validate.Validate) *Handler {
	return &Handler{s, v}
}
