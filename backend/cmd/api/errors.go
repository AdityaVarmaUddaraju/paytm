package main

import (
	"fmt"
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Error(
		err.Error(),
		"request_method", r.Method,
		"request_url", r.URL.String(),
	)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	data := envelope{
		"error": message,
	}

	err := app.writeJson(w, status, data, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("%s method is not allowed for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "resource not found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) invalidCreditialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) duplicateUsernameResponse(w http.ResponseWriter, r *http.Request) {
	message := "username already exists"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) invalidJWTTokenResponse(w http.ResponseWriter, r *http.Request, message string) {
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "user is not permitter to update"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to edit the user"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) duplicateAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "account already exists"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) accountMissingResponse(w http.ResponseWriter, r *http.Request) {
	message := "user does not have an account"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) insufficientBalanceResponse(w http.ResponseWriter, r *http.Request) {
	message := "insufficent balance"
	app.errorResponse(w, r, http.StatusConflict, message)
}
