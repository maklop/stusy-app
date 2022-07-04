package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	u "api/server/handlers/users"
	mw "api/server/middleware"
	"api/server/utils"
	"github.com/gorilla/mux"
)

type Route struct {
	Route   string
	Methods string
}

var r = mux.NewRouter()

func New() *mux.Router {

	r.HandleFunc("/", list).Methods("GET", "OPTIONS")
	r.HandleFunc("/users", u.Count()).Methods("GET", "OPTIONS")
	r.HandleFunc("/users/login", u.SignIn()).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/register", u.SignUp()).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/{id:[0-9]+}", mw.IsAuth(u.Info())).Methods("GET", "PUT", "OPTIONS")

	r.Use(mux.CORSMethodMiddleware(r))

	return r
}

func list(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		output := make(map[string]string, 15)
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			output["route"] = pathTemplate
		}
		methods, err := route.GetMethods()
		if err == nil {
			output["methods"] = strings.Join(methods, ",")
		}

		listJson(&w, output["route"], output["methods"])
		return nil
	})

	if err != nil {
		utils.WriteError(&w, http.StatusInternalServerError, utils.ErrorInternal)
		log.Println(err)
		return
	}

}

func listJson(w *http.ResponseWriter, path, methods string) {
	res, _ := json.MarshalIndent(struct {
		Path    string `json:"route"`
		Methods string `json:"methods"`
	}{
		Path:    path,
		Methods: methods,
	}, "", "	")

	(*w).Write(res)
}
