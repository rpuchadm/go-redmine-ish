package database

import (
	"database/sql"
	"fmt"
	"log"

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

	log.Println("Conexión a PostgreSQL establecida")

	// Crear la tabla projects si no existe
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS projects (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		identifier VARCHAR(255) UNIQUE NOT NULL,
		description TEXT,
		created_on TIMESTAMP DEFAULT NOW(),
		updated_on TIMESTAMP DEFAULT NOW(),
		parent_id INT,
		FOREIGN KEY (parent_id) REFERENCES projects(id) ON DELETE SET NULL
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("error al crear la tabla projects: %v", err)
	}

	log.Println("Tabla 'projects' creada o ya existente")

	return db, nil
}
