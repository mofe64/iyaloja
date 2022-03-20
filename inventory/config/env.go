package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var envLoaded = false

func loadEnv() {
	if envLoaded {
		return
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln("Error loading env file")
		}
		envLoaded = true
	}
}

func EnvMongoURI() string {
	loadEnv()
	return os.Getenv("MONGOURI")
}

func EnvHTTPPort() string {
	loadEnv()
	return os.Getenv("HTTP_PORT")
}

func EnvProfile() string {
	loadEnv()
	return os.Getenv("PROFILE")
}
