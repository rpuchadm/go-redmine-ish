package config

import (
	"fmt"
	"os"
)

type Config struct {
	AuthToken  string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() *Config {

	// Cargar la configuración desde las variables de entorno
	auth_token := os.Getenv("AUTH_TOKEN")
	if auth_token == "" {
		fmt.Println("ERROR AUTH_TOKEN no está definido")
		os.Exit(1)
	}

	db_host := os.Getenv("POSTGRES_SERVICE")
	if db_host == "" {
		fmt.Println("ERROR POSTGRES_SERVICE no está definido")
		os.Exit(1)
	}
	/*
		db_port := os.Getenv("DB_PORT")
		if db_port == "" {
			fmt.Println("ERROR DB_PORT no está definido")
			os.Exit(1)
		}*/
	db_user := os.Getenv("POSTGRES_USER")
	if db_user == "" {
		fmt.Println("ERROR POSTGRES_USER no está definido")
		os.Exit(1)
	}
	db_password := os.Getenv("POSTGRES_PASSWORD")
	if db_password == "" {
		fmt.Println("ERROR POSTGRES_PASSWORD no está definido")
		os.Exit(1)
	}
	db_name := os.Getenv("POSTGRES_DB")
	if db_name == "" {
		fmt.Println("ERROR POSTGRES_DB no está definido")
		os.Exit(1)
	}

	return &Config{
		AuthToken:  auth_token,
		DBHost:     db_host,
		DBPort:     "5432", // Puerto por defecto de PostgreSQL
		DBUser:     db_user,
		DBPassword: db_password,
		DBName:     db_name,
	}
}
