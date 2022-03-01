package application_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brice-74/golang-base-api/internal/api/application"
)

func TestWriteJsonContent(t *testing.T) {
	a := application.Application{}

	type payload struct {
		ID        int    `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	tests := []struct {
		Envelope    application.Envelope
		expected    string
		expectError error
	}{
		{
			Envelope: application.Envelope{"person": payload{1, "John", "Doe"}},
			expected: `{"person":{"id":1,"firstName":"John","lastName":"Doe"}}`},
		{
			Envelope: application.Envelope{"person": payload{2, "Kenzie", "Warner"}},
			expected: `{"person":{"id":2,"firstName":"Kenzie","lastName":"Warner"}}`,
		},
		{
			Envelope: application.Envelope{"person": payload{3, "Brice", "Butler"}},
			expected: `{"person":{"id":3,"firstName":"Brice","lastName":"Butler"}}`,
		},
		{
			Envelope:    application.Envelope{"channel": make(chan int)},
			expectError: errors.New("json: unsupported type: chan int"),
		},
		{
			Envelope:    application.Envelope{"inf": math.Inf(1)},
			expectError: errors.New("json: unsupported value: +Inf"),
		},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		err := a.WriteJSON(w, 200, tt.Envelope, nil)

		if tt.expectError != nil {
			if err == nil {
				t.Fatal("got nil error")
			}

			if err.Error() != tt.expectError.Error() {
				t.Errorf("got: %s, expect: %s", err.Error(), tt.expectError.Error())
			}
		} else {
			if err != nil {
				t.Fail()
			}

			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fail()
			}

			if strings.TrimSpace(string(body)) != strings.TrimSpace(tt.expected) {
				t.Errorf("got %s, expected %s", body, tt.expected)
			}
		}
	}
}

func TestWriteJsonStatus(t *testing.T) {
	a := application.Application{}

	tests := []int{200, 404, 500}

	for _, s := range tests {
		w := httptest.NewRecorder()

		err := a.WriteJSON(w, s, application.Envelope{}, nil)
		if err != nil {
			t.Fail()
		}

		ss := w.Result().StatusCode
		if ss != s {
			t.Errorf("got %d, expected %d", ss, s)
		}
	}
}

func TestWriteJsonHeaders(t *testing.T) {
	a := application.Application{}

	h := http.Header{}
	h.Add("Authorization", "Bearer 123")
	h.Add("Server", "Go")

	w := httptest.NewRecorder()

	err := a.WriteJSON(w, 200, application.Envelope{}, h)
	if err != nil {
		t.Fail()
	}

	for k, v := range h {
		value := w.Header().Get(k)
		if value != v[0] {
			t.Errorf("got %s for key %s, expected %s", value, k, v)
		}
	}
}

func TestReadJsonDecode(t *testing.T) {
	a := application.Application{}

	var input struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(`{"name":"John","age":24}`)))

	err := a.ReadJSON(w, r, &input)
	if err != nil {
		t.Fail()
	}

	if input.Name != "John" {
		t.Errorf("got %s, expected %s", input.Name, "John")
	}

	if input.Age != 24 {
		t.Errorf("got %d, expected %d", input.Age, 24)
	}
}

func TestReadJsonError(t *testing.T) {
	a := application.Application{}

	input := struct {
		Name string `json:"name"`
	}{Name: "John"}

	tests := []struct {
		json  string
		error string
	}{
		{json: `["name": 4]`, error: "body contains badly-formed JSON"},
		{json: `{"name":[]}`, error: "body contains incorrect JSON type"},
		{json: `{"age":24}`, error: "body contains unknown key"},
		{json: "", error: "body must not be empty"},
		{json: "{}{}", error: "body must only contain a single JSON value"},
	}

	for _, tt := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(tt.json)))

		err := a.ReadJSON(w, r, &input)
		if err == nil {
			t.Error("expected error, got nil")
		}

		if !strings.Contains(err.Error(), tt.error) {
			t.Errorf("got %s, expected %s", err.Error(), tt.error)
		}
	}
}
