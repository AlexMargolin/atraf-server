package account

import (
	"encoding/json"
	"net/http"

	"atraf-server/pkg/middleware"
	"atraf-server/pkg/token"
	"atraf-server/services/users"

	"atraf-server/pkg/rest"
	"atraf-server/pkg/validate"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type ActivateRequest struct {
	Code string `json:"code" validate:"required"`
}

type ActivateResponse struct {
	AccessToken string `json:"access_token"`
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

		account, err := handler.service.Register(request.Email, request.Password)
		if err != nil {
			rest.Error(w, http.StatusConflict)
			return
		}

		// Dependency(Users)
		// TODO this should in storage as a transaction? and then the mail won't be sent as well :)
		if u.NewUser(account.Id, users.UserFields{Email: account.Email}) != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		accessToken, err := token.NewAccessToken(token.AccessTokenCustomClaims{
			Active:    false,
			AccountId: account.Id,
		})
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusCreated, &RegisterResponse{
			accessToken,
		})
	}
}

func (handler *Handler) Activate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request ActivateRequest
		session := middleware.GetSessionContext(r)

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			rest.Error(w, http.StatusUnsupportedMediaType)
			return
		}

		if err := handler.validate.Struct(request); err != nil {
			rest.Error(w, http.StatusUnprocessableEntity)
			return
		}

		if err := handler.service.Activate(session.AccountId, request.Code); err != nil {
			rest.Error(w, http.StatusBadRequest)
			return
		}

		// Issue New Access Token
		accessToken, err := token.NewAccessToken(token.AccessTokenCustomClaims{
			Active:    true,
			AccountId: session.AccountId,
		})
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
			return
		}

		rest.Success(w, http.StatusOK, &ActivateResponse{
			AccessToken: accessToken,
		})
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

		// Issue Access Token
		accessToken, err := token.NewAccessToken(token.AccessTokenCustomClaims{
			Active:    account.Active,
			AccountId: account.Id,
		})
		if err != nil {
			rest.Error(w, http.StatusInternalServerError)
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
