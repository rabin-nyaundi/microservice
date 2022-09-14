package main

import (
	"log"
	"net/http"

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

	userCreated, err := app.models.User.Insert(user)
	if err != nil {
		log.Panic(err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated,
		jsonResponse{
			Error:   false,
			Success: true,
			Message: "user creation success",
			Data:    userCreated,
		})

	if err != nil {
		log.Panic(err)
		return
	}

}
