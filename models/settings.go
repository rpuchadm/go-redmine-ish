package models

import "database/sql"

/*
CREATE TABLE IF NOT EXISTS settings (
    id SERIAL PRIMARY KEY,              -- Identificador único de la configuración
    key VARCHAR(255) UNIQUE NOT NULL,   -- Clave de la configuración (por ejemplo, "theme", "language")
    value TEXT,                         -- Valor de la configuración
    created_at TIMESTAMP DEFAULT NOW(), -- Fecha de creación de la configuración
    updated_at TIMESTAMP DEFAULT NOW()  -- Fecha de última actualización de la configuración
);
*/

type Setting struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateSetting crea una nueva configuración
func CreateSetting(db *sql.DB, setting *Setting) (int, error) {
	query := `
	INSERT INTO settings (key, value)
	VALUES ($1, $2)
	RETURNING id`
	var id int
	err := db.QueryRow(query, setting.Key, setting.Value).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
