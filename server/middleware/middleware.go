package middleware

import (
	"api/config"
	"api/server/utils"

	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func IsAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		if r.Method == "OPTIONS" {
			return
		}

		header := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(header) != 2 || header[0] == "null" {
			utils.WriteError(&w, http.StatusBadRequest, utils.ErrorBadToken)
			return
		}

		if ok, _ := ValidateToken(header[1]); !ok {
			utils.WriteError(&w, http.StatusUnauthorized, utils.ErrorBadToken)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func CreateToken(issuer string) (string, int64, error) {
	expiresAt := time.Now().Add(time.Hour * 1).Unix()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: expiresAt,
	})

	token, err := claims.SignedString([]byte(config.Secret))
	if err != nil {
		return "", 0, err
	}
	return token, expiresAt, nil
}

func ValidateToken(tokenString string) (bool, error) {
	token, err := parseToken(tokenString)
	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, nil
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Secret), nil
	})

	return token, err
}
