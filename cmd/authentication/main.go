package main

import (
	"log"

	"quotes/config"
)

type Config struct {
	DbHost       string `env:"MYSQL_HOST"`     // Database host name
	DbPort       int    `env:"MYSQL_PORT"`     // Database port
	DbUser       string `env:"MYSQL_USER"`     // Database username
	DbPass       string `env:"MYSQL_PASS"`     // Database user password
	DbName       string `env:"MYSQL_NAME"`     // Database DB Name
	ServerPort   int    `env:"SERVER_PORT"`    // Server Port
	AuthATSecret string `env:"AUTH_AT_SECRET"` // Authentication Access Token Signing Secret
	AuthRTSecret string `env:"AUTH_RT_SECRET"` // Authentication Refresh Token Signing Secret
}

func main() {
	// Default Configuration
	cfg := &Config{
		DbHost:       "localhost",
		DbPort:       3306,
		DbUser:       "root",
		DbPass:       "",
		DbName:       "at_db",
		ServerPort:   8080,
		AuthATSecret: "abcd",
		AuthRTSecret: "dcba",
	}

	if err := config.Marshal(cfg); err != nil {
		log.Fatal(err)
	}
}
