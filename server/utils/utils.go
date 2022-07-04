package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	ErrorInternal       = "Internal server error. Try again later."
	ErrorBadJSON        = "You have supplied an invalid JSON body."
	ErrorBadCredentials = "Invalid credentials. Invalid email or password."
	ErrorUserExist      = "Email is already in use."
	ErrorEmailFormat    = "Invalid email format."
	ErrorShortPass      = "Password is too short (should be >= 8)."
	ErrorBadToken       = "Token is invalid or expired."
	ErrorAccess         = "You don't have permissions to make this request."
	ErrorNameFormat     = "Name should only consist of alphabetical characters."
	ErrorNotFound       = "Resource not found."
)

func WriteError(w *http.ResponseWriter, status int, err string) {
	res, _ := json.MarshalIndent(struct {
		Message string `json:"error_message"`
	}{
		Message: err,
	}, "", "	")

	(*w).WriteHeader(status)
	(*w).Write(res)
}

func ReadJSON(w *http.ResponseWriter, r *http.Request) (map[string]string, error) {
	var data map[string]string

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, ErrorInternal)
		return nil, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		WriteError(w, http.StatusBadRequest, ErrorBadJSON)
		return nil, err
	}

	return data, nil
}
