package users

import "net/http"

func Count() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Simply return current users count
	}
}
