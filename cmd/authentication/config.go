package main

import (
	"log"

	"quotes/config"
)

type Config struct {
	DbHost       string `env:"MYSQL_HOST"`
	DbPort       int    `env:"MYSQL_PORT"`
	DbUser       string `env:"MYSQL_USER"`
	DbPass       string `env:"MYSQL_PASS"`
	DbName       string `env:"MYSQL_NAME"`
	ServerPort   int    `env:"SERVER_PORT"`
	AuthATSecret string `env:"AUTH_AT_SECRET"`
	AuthRTSecret string `env:"AUTH_RT_SECRET"`
}

func main() {
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
