package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/signup", app.userRegisterHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/signin", app.userSignInHandler)

	router.HandlerFunc(http.MethodGet, "/v1/users", app.listUsersHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/:id", app.authenticate(app.updateUserHandler))

	router.HandlerFunc(http.MethodPost, "/v1/accounts/create", app.authenticate(app.createAccountHandler))
	router.HandlerFunc(http.MethodPost, "/v1/accounts/add", app.authenticate(app.addMoneyHandler))
	router.HandlerFunc(http.MethodPost, "/v1/accounts/transfer", app.authenticate(app.transferMoneyHandler))

	return app.recoverPanic(app.enableCORS(router))
}
