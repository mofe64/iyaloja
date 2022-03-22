package config

import (
	"github.com/joho/godotenv"
	"github.com/mofe64/iyaloja/inventory/util"
	"os"
)

var envLoaded = false

func loadEnv() {
	if envLoaded {
		return
	} else {
		err := godotenv.Load()
		if err != nil {
			util.ApplicationLog.Fatalf("Error loading env file %v\n", err)
		}
		envLoaded = true
	}
}

func EnvMongoURI() string {
	loadEnv()
	return os.Getenv("MONGOURI")
}

func EnvDatabaseName() string {
	return os.Getenv("DATABASE_NAME")
}

func EnvHTTPPort() string {
	loadEnv()
	return os.Getenv("HTTP_PORT")
}

func EnvProfile() string {
	loadEnv()
	return os.Getenv("PROFILE")
}
