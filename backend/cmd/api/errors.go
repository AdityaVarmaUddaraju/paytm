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