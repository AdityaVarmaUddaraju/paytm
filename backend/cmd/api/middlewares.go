package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/AdityaVarmaUddaraju/paytm/internal/tokens"
)

func (app *application) authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Vary", "Authorization")

		// check if authentication header is provided
		authenticationHeader := r.Header.Get("Authorization")

		if authenticationHeader == "" {
			app.invalidJWTTokenResponse(w, r, authenticationHeader)
			return
		}

		authParts := strings.Split(authenticationHeader, " ")
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			app.invalidJWTTokenResponse(w, r, authenticationHeader)
			return
		}
		// verify jwt token in authentication header
		token := authParts[1]

		username, err := tokens.VerifyToken(token, app.cfg.jwtSecretKey)

		if err != nil {
			switch {
			case errors.Is(err, tokens.ErrInvalidJWTToken):
				app.invalidJWTTokenResponse(w, r, err.Error())
				return
			default:
				app.serverErrorResponse(w, r, err)
				return
			}

		}

		// get the user from username in jwt
		user, err := app.models.Users.GetByUsername(username)

		if err != nil {
			app.invalidJWTTokenResponse(w, r, err.Error())
			return
		}
		// set the user to request context
		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)

	})
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Vary", "Origin")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range app.cfg.cors.trustedOrigins {
				if origin == app.cfg.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						w.WriteHeader(http.StatusOK)
						return
					}

					break
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
