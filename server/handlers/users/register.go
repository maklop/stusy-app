package users

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"api/server/models"
	"api/server/utils"

	"golang.org/x/crypto/bcrypt"
	_ "gorm.io/driver/mysql"
)

func SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			return
		}

		data, err := utils.ReadJSON(&w, r)
		if err != nil {
			return
		}

		if len(data["password"]) < 8 {
			utils.WriteError(&w, http.StatusBadRequest, utils.ErrorShortPass)
			return
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 12)

		re := regexp.MustCompile(`^[a-zA-Z\d.\-]+@[a-zA-Z\d\-]+\.[a-zA-Z.]`).MatchString
		if !re(data["email"]) {
			utils.WriteError(&w, http.StatusBadRequest, utils.ErrorEmailFormat)
			return
		}

		user := models.User{
			Email:    strings.ToLower(data["email"]),
			Password: password,
		}

		var temp models.User
		models.DB.Where("email = ?", user.Email).First(&temp)
		if temp.ID != 0 {
			utils.WriteError(&w, http.StatusForbidden, utils.ErrorUserExist)
			return
		}

		models.DB.Create(&user)

		registerMessage(&w, &user.ID)
	}
}

func registerMessage(w *http.ResponseWriter, id *uint) {
	res, _ := json.MarshalIndent(struct {
		UserID uint   `json:"user_id"`
		Href   string `json:"href"`
	}{
		UserID: *id,
		Href:   "https://api.studentsystem.xyz/profile/" + strconv.Itoa(int(*id)),
	}, "", "	")

	(*w).WriteHeader(http.StatusCreated)
	(*w).Write(res)
}
