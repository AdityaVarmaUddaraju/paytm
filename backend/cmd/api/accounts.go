package main

import (
	"errors"
	"net/http"

	"github.com/AdityaVarmaUddaraju/paytm/internal/data"
	"github.com/AdityaVarmaUddaraju/paytm/internal/validator"
)

func (app *application) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	err := app.models.Accounts.CreateAccount(user.ID)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateAccount):
			app.duplicateAccountResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	data := envelope{
		"message": "account created successfully",
	}

	app.writeJson(w, http.StatusCreated, data, nil)
}

func (app *application) addMoneyHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var input struct {
		Amount int `json:"amount"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.Amount > 0, "amount", "should be greater than or equal to 0")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Accounts.AddMoney(user.ID, input.Amount)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoAccount):
			app.accountMissingResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	data := envelope{
		"message": "money added successfully",
	}

	app.writeJson(w, http.StatusOK, data, nil)
}

func (app *application) transferMoneyHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var input struct {
		UserID int64 `json:"user_id"`
		Amount int `json:"amount"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ok, err := app.models.Accounts.CheckIfUserExists(input.UserID)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoAccount):
			app.accountMissingResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	if !ok {
		app.accountMissingResponse(w, r)
		return
	}

	v := validator.New()

	v.Check(input.Amount > 0, "amount", "amount should be greater than 0")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Accounts.TransferMoney(user.ID, input.UserID, input.Amount)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrInsuffientBalance):
			app.insufficientBalanceResponse(w, r)
			return
		case errors.Is(err, data.ErrNoAccount):
			app.accountMissingResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}

	}

	data := envelope{
		"message": "amount transfered successfully",
	}

	err = app.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
