package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/rabin-nyaundi/authentication-service/internal/data"
)

// createUserHandeler adds a user to the database and a tokn to the tokens table
func (app *application) createUserHandeler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)

	if err != nil {
		log.Panic(err)
		return
	}

	user := &data.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Active:    false,
	}

	err = user.Password.Set(input.Password)

	if err != nil {
		log.Panic(err)
		return
	}

	err = app.models.User.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, data.DuplicateEmail):
			app.writeJSON(w, http.StatusBadRequest, JSONResponse{
				Error:   true,
				Message: "user wit email already exist",
			})
		default:
			log.Panic(err)
		}
		return
	}

	duration := 1 * 24 * time.Hour
	token, err := app.models.Token.New(user.ID, duration, data.ScopeActivation)

	if err != nil {
		log.Panic(err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated,
		JSONResponse{
			Success: true,
			Message: "user creation success",
			Data:    token,
		})

	if err != nil {
		log.Panic(err)
		return
	}

}

// authenticateHandler checks the given user details against the databse if the match and return json response
func (app *application) authenticateHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.JSONEror(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)

	log.Println(user)

	if err != nil {
		app.JSONEror(w, err, http.StatusBadRequest)
		return
	}

	valid, err := user.Password.MatchesPassword(input.Password)

	if err != nil || !valid {
		app.JSONEror(w, err, http.StatusUnauthorized)
		return
	}

	app.writeJSON(w, http.StatusAccepted, JSONResponse{
		Success: true,
		Message: "user authentication success",
		Data:    user,
	})

}

// fetchUserHandler returns a singl user from database
func (app *application) fetchUserHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParams(w, r)

	if err != nil {
		log.Panic(err)
		return
	}

	user, err := app.models.User.GetOneUser(int(id))

	if err != nil {
		log.Panic(err)
		return
	}

	app.writeJSON(w, http.StatusOK, JSONResponse{
		Success: true,
		Message: "success",
		Data:    user,
	})
}
