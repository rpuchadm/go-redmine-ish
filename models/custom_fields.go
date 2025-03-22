package models

import "database/sql"

/*
CREATE TABLE IF NOT EXISTS custom_fields (
    id SERIAL PRIMARY KEY,              -- Identificador único del campo personalizado
    name VARCHAR(255) NOT NULL,         -- Nombre del campo personalizado
    field_type VARCHAR(50) NOT NULL,    -- Tipo de campo (por ejemplo, "text", "number", "date")
    default_value TEXT,                 -- Valor por defecto del campo
    is_required BOOLEAN DEFAULT FALSE,  -- Indica si el campo es obligatorio
    created_at TIMESTAMP DEFAULT NOW(), -- Fecha de creación del campo
    updated_at TIMESTAMP DEFAULT NOW()  -- Fecha de última actualización del campo
);
*/

type CustomField struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	FieldType    string `json:"field_type"`
	DefaultValue string `json:"default_value"`
	IsRequired   bool   `json:"is_required"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateCustomField crea un nuevo campo personalizado
func CreateCustomField(db *sql.DB, customField *CustomField) (int, error) {
	query := `
	INSERT INTO custom_fields (name, field_type, default_value, is_required)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var id int
	err := db.QueryRow(query, customField.Name, customField.FieldType, customField.DefaultValue, customField.IsRequired).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetCustomFieldByID(db *sql.DB, id int) (*CustomField, error) {
	query := `
	SELECT id, name, field_type, default_value, is_required, created_at, updated_at
	FROM custom_fields
	WHERE id = $1`
	customField := &CustomField{}
	err := db.QueryRow(query, id).Scan(&customField.ID, &customField.Name, &customField.FieldType, &customField.DefaultValue, &customField.IsRequired, &customField.CreatedAt, &customField.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return customField, nil
}

func GetCustomFields(db *sql.DB) ([]CustomField, error) {
	query := `
	SELECT id, name, field_type, default_value, is_required, created_at, updated_at
	FROM custom_fields`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customFields := []CustomField{}
	for rows.Next() {
		var customField CustomField
		if err := rows.Scan(&customField.ID, &customField.Name, &customField.FieldType, &customField.DefaultValue, &customField.IsRequired, &customField.CreatedAt, &customField.UpdatedAt); err != nil {
			return nil, err
		}
		customFields = append(customFields, customField)
	}

	return customFields, nil
}

func UpdateCustomField(db *sql.DB, customField *CustomField) error {
	query := `
	UPDATE custom_fields
	SET name = $1, field_type = $2, default_value = $3, is_required = $4, updated_at = NOW()
	WHERE id = $5`
	_, err := db.Exec(query, customField.Name, customField.FieldType, customField.DefaultValue, customField.IsRequired, customField.ID)
	return err
}

func DeleteCustomField(db *sql.DB, id int) error {
	query := `DELETE FROM custom_fields WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func CountCustomFields(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM custom_fields`
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CreateCustomFieldTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS custom_fields (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		field_type VARCHAR(50) NOT NULL,
		default_value TEXT,
		is_required BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	)`
	_, err := db.Exec(query)
	return err
}

func DropCustomFieldTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS custom_fields`
	_, err := db.Exec(query)
	return err
}
