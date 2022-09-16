package main

import (
	"log"
	"net/http"
	"time"

	"github.com/rabin-nyaundi/authentication-service/cmd/internal/data"
)

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

	_, err = app.models.User.Insert(user)
	if err != nil {
		log.Panic(err)
		return
	}

	duration := 1 * 24 * time.Hour
	token, err := app.models.Token.New(user.ID, duration, data.ScopeActivation)

	if err != nil {
		log.Panic(err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated,
		jsonResponse{
			Error:   false,
			Success: true,
			Message: "user creation success",
			Data:    token,
		})

	if err != nil {
		log.Panic(err)
		return
	}

}
