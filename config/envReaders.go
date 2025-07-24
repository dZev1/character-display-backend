package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

func ReadConnStrEnv() (string, error) {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		return "", err
	}

	envPGUser := envFile["PGUSER"]
	envPGPassword := envFile["PGPASSWORD"]

	connStr := fmt.Sprintf("postgresql://%v:%v@ep-plain-snowflake-acxitlrh-pooler.sa-east-1.aws.neon.tech/myDB?sslmode=require", envPGUser, envPGPassword)
	return connStr, nil
}

func ReadPortEnv() (string, error) {
	envFile, err := godotenv.Read(".env")
	if err != nil {
		return "", err
	}

	return envFile["PORT"], nil
}
