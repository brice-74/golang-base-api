package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// LogError uses the normal logger and adds details about current request
func (app *Application) LogError(r *http.Request, err error) {
	app.Logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

// WriteJSON is an helper simplifying writing JSON to an HTTP response.
func (app *Application) WriteJSON(w http.ResponseWriter, status int, data Envelope, headers http.Header) error {

	for key, value := range headers {
		w.Header()[key] = value
	}

	if data != nil {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(status)

	if data != nil {
		js, err := json.Marshal(data)
		if err != nil {
			return err
		}

		js = append(js, '\n')

		_, err = w.Write(js)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadJSON is an helper for simplifying JSON reading from a HTTP body and error handling during decoding.
func (app *Application) ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Limit the size to the request body (1mb)
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Prevent unwanted fields in the body
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "application: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}
