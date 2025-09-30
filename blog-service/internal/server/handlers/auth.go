package handlers

import (
	pb "blog-service/internal/grpc/protobuf"
	"blog-service/internal/server/models"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type AuthHandler struct {
	AuthClient pb.AuthServiceClient
}

const (
	registerPath      = "/auth/register"
	registerPathSlash = "/auth/register/"
	loginPath         = "/auth/login"
	loginPathSlash    = "/auth/login/"
)

const bearerSchema = "Bearer "

/*
AuthMiddleware attempts to extract and validate a jwt token from the authorization header.
If it succeeds, it will send the claims to the handlers via the request context.
If it does not, it will send a nil pointer.
This way, the auth middleware can be used on handlers that handle both protected
and unprotected routes, authorization being established using other helper functions.
*/
func AuthMiddleware(h http.Handler, authClient pb.AuthServiceClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userClaims *models.UserClaims

		defer func() {
			ctx := context.WithValue(r.Context(), models.ClaimsKey, userClaims)
			h.ServeHTTP(w, r.WithContext(ctx))
		}()

		headerAccess := r.Header.Get("Authorization")
		if len(headerAccess) <= len(bearerSchema) ||
			!strings.HasPrefix(headerAccess, bearerSchema) {
			return
		}

		tokenString := headerAccess[len(bearerSchema):]
		userClaims = models.VerifyAuthToken(r.Context(), authClient, tokenString)
	})
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && (r.URL.Path == registerPath || r.URL.Path == registerPathSlash):
		h.AuthRegister(w, r)
		return
	case r.Method == http.MethodPost && (r.URL.Path == loginPath || r.URL.Path == loginPathSlash):
		h.AuthLogin(w, r)
		return
	}

	http.NotFound(w, r)
}

func (h *AuthHandler) AuthLogin(w http.ResponseWriter, r *http.Request) {
	var loginDTO models.UserLoginDTO

	err := json.NewDecoder(r.Body).Decode(&loginDTO)
	if err != nil {
		http.Error(w, "invalid request sent", http.StatusBadRequest)
		return
	}

	var ParamErr *models.ParamError
	var InvalidLoginErr *models.InvalidLoginError

	token, err := models.LoginUser(r.Context(), h.AuthClient, &loginDTO)

	switch {
	case errors.As(err, &ParamErr):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.As(err, &InvalidLoginErr):
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(models.UserLoginResponseDTO{Token: token})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *AuthHandler) AuthRegister(w http.ResponseWriter, r *http.Request) {
	var registerDto models.UserRegisterDTO

	err := json.NewDecoder(r.Body).Decode(&registerDto)
	if err != nil {
		http.Error(w, "invalid request sent", http.StatusBadRequest)
		return
	}

	var ParamErr *models.ParamError
	var EmailOrUserTakenErr *models.EmailOrUserTakenError

	err = models.RegisterUser(r.Context(), h.AuthClient, &registerDto)

	if err != nil {
		if errors.As(err, &ParamErr) || errors.As(err, &EmailOrUserTakenErr) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}
}
