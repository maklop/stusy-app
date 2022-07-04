package users

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"api/server/models"
	"api/server/utils"
	"github.com/gorilla/mux"
)

func Info() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user models.User
		var userData models.UserData

		models.DB.Where("token = ?", token[1]).First(&user)
		if user.ID != uint(id) {
			utils.WriteError(&w, http.StatusForbidden, utils.ErrorAccess)
			return
		}

		if r.Method == "GET" {
			models.DB.Where("user_id = ?", uint(id)).First(&userData)
			if userData.UserID == 0 {
				utils.WriteError(&w, http.StatusNotFound, utils.ErrorNotFound)
				return
			}

			GetResponse(&w, &userData.FirstName, &userData.LastName)
			return
		}

		data, err := utils.ReadJSON(&w, r)
		if err != nil {
			return
		}

		re := regexp.MustCompile(`^\p{L}+$`).MatchString
		if !re(data["first_name"]) || !re(data["last_name"]) {
			utils.WriteError(&w, http.StatusBadRequest, utils.ErrorNameFormat)
			return
		}

		models.DB.Where("user_id = ?", uint(id)).First(&userData)

		userData.LastName = data["last_name"]
		userData.FirstName = data["first_name"]

		if userData.ID == 0 {
			userData.UserID = user.ID
			models.DB.Create(&userData)
			w.WriteHeader(http.StatusCreated)
			return
		}

		w.WriteHeader(http.StatusForbidden)
	}
}

func GetResponse(w *http.ResponseWriter, first, last *string) {
	res, _ := json.MarshalIndent(struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{
		FirstName: *first,
		LastName:  *last,
	}, "", "	")

	(*w).WriteHeader(http.StatusOK)
	(*w).Write(res)
}
