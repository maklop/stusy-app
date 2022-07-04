package config

import (
	"fmt"
	"os"
)

var (
	Secret string
	Port   string
	DSN    string
)

func Init() error {
	DSN = fmt.Sprintf("%s:%s@tcp(stusy-db)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	Port = ":" + os.Getenv("PORT")
	Secret = os.Getenv("SECRET")

	if Port == "" || Secret == "" {
		return fmt.Errorf("some of environment variables are blank or missing (make sure to set them right)")
	}

	return nil
}
