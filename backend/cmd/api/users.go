package main

import (
	"errors"
	"net/http"

	"github.com/AdityaVarmaUddaraju/paytm/internal/data"
	"github.com/AdityaVarmaUddaraju/paytm/internal/tokens"
	"github.com/AdityaVarmaUddaraju/paytm/internal/validator"
)

func (app *application) userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserName  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		UserName:  input.UserName,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	err = user.Password.Set(input.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// validate user input
	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateUsername):
			app.duplicateUsernameResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}

	}

	data := envelope{
		"user": user,
	}

	err = app.writeJson(w, http.StatusCreated, data, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) userSignInHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validate input
	v := validator.New()

	v.ValidateEmpty(input.Username, "username")
	data.ValidatePassword(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// check if user exists with given username if not raise invalid credential
	user, err := app.models.Users.GetByUsername(input.Username)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCreditialsResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// check if password matches with hashed password if not raise invalid credentials
	match, err := user.Password.Match(input.Password)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCreditialsResponse(w, r)
		return
	}

	// if user is valid send the jwt token in response
	token, err := tokens.CreateToken(app.cfg.jwtSecretKey, user.UserName)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//send the token to the user
	data := envelope{
		"token": token,
	}
	app.writeJson(w, http.StatusOK, data, nil)
}

func (app *application) listUsersHandler(w http.ResponseWriter, r *http.Request) {

	searchTerm := r.URL.Query().Get("search")

	users, err := app.models.Users.GetUsers(searchTerm)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"users": users,
	}

	err = app.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	id, err := app.readIdParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	if id != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	user.FirstName = input.Firstname
	user.LastName = input.Lastname
	user.Password.Set(input.Password)

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.UpdateUser(user)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.editConflictResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

}

