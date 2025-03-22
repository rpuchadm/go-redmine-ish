package database

import (
	"database/sql"
	"fmt"

	"go-redmine-ish/config"

	_ "github.com/lib/pq" // Driver de PostgreSQL
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	// Cadena de conexión a PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Conectar a la base de datos
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos: %v", err)
	}

	// Verificar la conexión
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error al verificar la conexión: %v", err)
	}

	// log.Println("Conexión a PostgreSQL establecida")

	return db, nil
}
