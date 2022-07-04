package users

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"api/server/middleware"
	"api/server/models"
	"api/server/utils"
	"golang.org/x/crypto/bcrypt"
	_ "gorm.io/driver/mysql"
)

func SignIn() http.HandlerFunc {
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

		var user models.User

		models.DB.Where("email = ?", strings.ToLower(data["email"])).First(&user)

		if user.ID == 0 {
			utils.WriteError(&w, http.StatusForbidden, utils.ErrorBadCredentials)
			return
		}

		if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
			utils.WriteError(&w, http.StatusForbidden, utils.ErrorBadCredentials)
			return
		}

		if len(user.Token) > 0 {
			if ok, _ := middleware.ValidateToken(user.Token); ok {
				loginMessage(&w, &user.Token, &user.ID, &user.ExpiresAt)
				return
			}
		}

		user.Token, user.ExpiresAt, err = middleware.CreateToken(strconv.Itoa(int(user.ID)))
		if err != nil {
			log.Println(err)
			utils.WriteError(&w, http.StatusInternalServerError, utils.ErrorInternal)
			return
		}

		models.DB.Save(&user)

		loginMessage(&w, &user.Token, &user.ID, &user.ExpiresAt)
	}
}

func loginMessage(w *http.ResponseWriter, token *string, id *uint, expiresAt *int64) {
	res, _ := json.MarshalIndent(struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		UserID      uint   `json:"user_id"`
	}{
		TokenType:   "bearer",
		AccessToken: *token,
		ExpiresIn:   *expiresAt - time.Now().Unix(),
		UserID:      *id,
	}, "", "	")

	(*w).WriteHeader(http.StatusOK)
	(*w).Write(res)
}
