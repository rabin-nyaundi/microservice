package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// JSONResponse structure for holding json response
type JSONResponse struct {
	Error   bool        `json:"error,omitempty"`
	Success bool        `json:"success,omitempty"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// RequestPayload holds request payload
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

// AuthPayload holds authentication request payload
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *application) Broker(w http.ResponseWriter, r *http.Request) {
	payload := JSONResponse{
		Error:   false,
		Message: "Successful",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}

func (app *application) submitRequestHandler(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	// authentication credentials
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.JSONEror(w, err, http.StatusBadRequest)
		return
	}

	// get user id from the request
	// userId, err := app.readIDParams(w, r)
	// if err != nil {
	// 	app.JSONEror(w, errors.New("invalid id"))
	// 	return
	// }

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	// case "getuser":
	// 	app.GetUser(w, userId)
	default:
		app.JSONEror(w, errors.New("Failed"))
	}
}

func (app *application) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, err := json.MarshalIndent(a, "", "\t")

	if err != nil {
		app.JSONEror(w, errors.New("error at mashal indent"))
		return
	}

	// call auth service
	request, err := http.NewRequest("POST", "http://authentication-service/v1/users/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.JSONEror(w, err, http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.JSONEror(w, err)
		return
	}

	defer response.Body.Close()

	// make sure we get correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.JSONEror(w, errors.New("invalid credentials"))
		return

	} else if response.StatusCode != http.StatusOK {
		app.JSONEror(w, errors.New("error calling auth service"))
		return
	}

	var jsonFromService JSONResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		log.Println("error at decode here")
		app.JSONEror(w, errors.New("error decoding request body"))
		return
	}

	if jsonFromService.Error {
		app.JSONEror(w, err, http.StatusForbidden)
	}

	var payload JSONResponse
	payload.Success = true
	payload.Message = "login succsessful"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}

// func (app *application) GetUser(w http.ResponseWriter, id int64) {

// 	request, err := http.NewRequest("GET", fmt.Sprintf("http://authentication-service/v1/users/%d", id))
// 	if err != nil {
// 		app.JSONEror(w, errors.New("invalid id"))
// 	}

// 	client := &http.Client{}
// 	response, err := client.Do(request)

// 	if err != nil {
// 		app.JSONEror(w, errors.New("invalid id"))
// 	}

// 	defer response.Body.Close()

// 	if response.StatusCode == http.StatusUnauthorized {
// 		app.JSONEror(w, errors.New("invalid credentials"))
// 		return

// 	} else if response.StatusCode != http.StatusOK {
// 		app.JSONEror(w, errors.New("error calling auth service"))
// 		return
// 	}

// 	var jsonFromService JSONResponse

// 	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
// 	if err != nil {
// 		log.Println("error at decode here")
// 		app.JSONEror(w, errors.New("error decoding request body"))
// 		return
// 	}

// 	if jsonFromService.Error {
// 		app.JSONEror(w, err, http.StatusForbidden)
// 	}

// 	var payload JSONResponse
// 	payload.Success = true
// 	payload.Message = "user fetch succsessful"
// 	payload.Data = jsonFromService.Data

// 	app.writeJSON(w, http.StatusAccepted, payload)

// }
