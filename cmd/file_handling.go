package main

import (
	"fmt"

	"github.com/joho/godotenv"
)

func connStrEnv() (string, error) {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		return "", err
	}

	envPGUser := envFile["PGUSER"]
	envPGPassword := envFile["PGPASSWORD"]

	connStr := fmt.Sprintf("postgresql://%v:%v@ep-plain-snowflake-acxitlrh-pooler.sa-east-1.aws.neon.tech/myDB?sslmode=require", envPGUser, envPGPassword)
	return connStr, nil
}