package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// JSONResponse structure holds response sent to a clint

// readIDParams returns the id from the get request
func (app *application) readIDParams(w http.ResponseWriter, r *http.Request) (int64, error) {
	userID := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(userID, 10, 12)

	if err != nil {
		log.Panic(err)
		return 0, err
	}

	fmt.Println(r.URL)
	return id, nil

}

// writeJSON converts data provided to jon format
func (app *application) writeJSON(w http.ResponseWriter, status int, data any) error {
	jsonObject, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		log.Panic(err)
		return err
	}

	jsonObject = append(jsonObject, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonObject)
	return nil
}

// readJSON reads the json object into struct
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1_048_576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalerror *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field == "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request too large":
			return fmt.Errorf("body must not be larger that %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalerror):
			panic(err)

		default:
			fmt.Println("Heeey")
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return fmt.Errorf("body must contain a single JSON object")
	}

	return nil
}

func (app *application) JSONEror(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)

}
