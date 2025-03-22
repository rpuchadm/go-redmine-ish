package models

import "database/sql"

/*
CREATE TABLE IF NOT EXISTS custom_field_values (
    id SERIAL PRIMARY KEY,              -- Identificador único del valor
    custom_field_id INT NOT NULL,       -- ID del campo personalizado
    entity_type VARCHAR(50) NOT NULL,   -- Tipo de entidad (por ejemplo, "project", "issue", "user")
    entity_id INT NOT NULL,             -- ID de la entidad asociada
    value TEXT,                         -- Valor del campo personalizado
    created_at TIMESTAMP DEFAULT NOW(), -- Fecha de creación del valor
    updated_at TIMESTAMP DEFAULT NOW(), -- Fecha de última actualización del valor
    FOREIGN KEY (custom_field_id) REFERENCES custom_fields(id) ON DELETE CASCADE
);
*/

// CustomFieldValue representa el valor de un campo personalizado
type CustomFieldValue struct {
	ID            int    `json:"id"`
	CustomFieldID int    `json:"custom_field_id"`
	EntityType    string `json:"entity_type"`
	EntityID      int    `json:"entity_id"`
	Value         string `json:"value"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CreateCustomFieldValue crea un nuevo valor de campo personalizado
func CreateCustomFieldValue(db *sql.DB, customFieldValue *CustomFieldValue) (int, error) {
	query := `
	INSERT INTO custom_field_values (custom_field_id, entity_type, entity_id, value)
	VALUES ($1, $2, $3, $4)
	RETURNING id`
	var id int
	err := db.QueryRow(query, customFieldValue.CustomFieldID, customFieldValue.EntityType, customFieldValue.EntityID, customFieldValue.Value).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetCustomFieldValuesByEntity obtiene todos los valores de campo personalizado de una entidad
func GetCustomFieldValuesByEntity(db *sql.DB, entityType string, entityID int) ([]CustomFieldValue, error) {
	query := `
	SELECT id, custom_field_id, entity_type, entity_id, value, created_at, updated_at
	FROM custom_field_values
	WHERE entity_type = $1 AND entity_id = $2`
	rows, err := db.Query(query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customFieldValues := []CustomFieldValue{}
	for rows.Next() {
		var customFieldValue CustomFieldValue
		if err := rows.Scan(&customFieldValue.ID, &customFieldValue.CustomFieldID, &customFieldValue.EntityType, &customFieldValue.EntityID, &customFieldValue.Value, &customFieldValue.CreatedAt, &customFieldValue.UpdatedAt); err != nil {
			return nil, err
		}
		customFieldValues = append(customFieldValues, customFieldValue)
	}

	return customFieldValues, nil
}

// GetCustomFieldValueByID obtiene un valor de campo personalizado por su ID
func GetCustomFieldValueByID(db *sql.DB, id int) (*CustomFieldValue, error) {
	query := `
	SELECT id, custom_field_id, entity_type, entity_id, value, created_at, updated_at
	FROM custom_field_values
	WHERE id = $1`
	customFieldValue := &CustomFieldValue{}
	err := db.QueryRow(query, id).Scan(&customFieldValue.ID, &customFieldValue.CustomFieldID, &customFieldValue.EntityType, &customFieldValue.EntityID, &customFieldValue.Value, &customFieldValue.CreatedAt, &customFieldValue.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return customFieldValue, nil
}

// UpdateCustomFieldValue actualiza un valor de campo personalizado
func UpdateCustomFieldValue(db *sql.DB, customFieldValue *CustomFieldValue) error {
	query := `
	UPDATE custom_field_values
	SET value = $1, updated_at = NOW()
	WHERE id = $2`
	_, err := db.Exec(query, customFieldValue.Value, customFieldValue.ID)
	return err
}

// DeleteCustomFieldValue elimina un valor de campo personalizado
func DeleteCustomFieldValue(db *sql.DB, id int) error {
	query := `DELETE FROM custom_field_values WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

// GetCustomFieldValuesByCustomFieldID obtiene todos los valores de campo personalizado de un campo personalizado
func GetCustomFieldValuesByCustomFieldID(db *sql.DB, customFieldID int) ([]CustomFieldValue, error) {
	query := `
	SELECT id, custom_field_id, entity_type, entity_id, value, created_at, updated_at
	FROM custom_field_values
	WHERE custom_field_id = $1`
	rows, err := db.Query(query, customFieldID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customFieldValues := []CustomFieldValue{}
	for rows.Next() {
		var customFieldValue CustomFieldValue
		if err := rows.Scan(&customFieldValue.ID, &customFieldValue.CustomFieldID, &customFieldValue.EntityType, &customFieldValue.EntityID, &customFieldValue.Value, &customFieldValue.CreatedAt, &customFieldValue.UpdatedAt); err != nil {
			return nil, err
		}
		customFieldValues = append(customFieldValues, customFieldValue)
	}

	return customFieldValues, nil
}

/*
// DeleteCustomFieldValuesByEntity elimina todos los valores de campo personalizado de una entidad
func DeleteCustomFieldValuesByEntity(db *sql.DB, entityType string, entityID int) error {
	query := `DELETE FROM custom_field_values WHERE entity_type = $1 AND entity_id = $2`
	_, err := db.Exec(query, entityType, entityID)
	return err
}

// DeleteCustomFieldValuesByCustomFieldID elimina todos los valores de campo personalizado de un campo personalizado
func DeleteCustomFieldValuesByCustomFieldID(db *sql.DB, customFieldID int) error {
	query := `DELETE FROM custom_field_values WHERE custom_field_id = $1`
	_, err := db.Exec(query, customFieldID)
	return err
}
*/

// DeleteCustomFieldValuesByCustomFieldIDAndEntity elimina todos los valores de campo personalizado de un campo personalizado y una entidad
func DeleteCustomFieldValuesByCustomFieldIDAndEntity(db *sql.DB, customFieldID int, entityType string, entityID int) error {
	query := `DELETE FROM custom_field_values WHERE custom_field_id = $1 AND entity_type = $2 AND entity_id = $3`
	_, err := db.Exec(query, customFieldID, entityType, entityID)
	return err
}

func CreateCustomFieldValuesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS custom_field_values (
		id SERIAL PRIMARY KEY,
		custom_field_id INT NOT NULL,
		entity_type VARCHAR(50) NOT NULL,
		entity_id INT NOT NULL,
		value TEXT,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		FOREIGN KEY (custom_field_id) REFERENCES custom_fields(id) ON DELETE CASCADE
	)`
	_, err := db.Exec(query)
	return err
}

func DropCustomFieldValuesTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS custom_field_values`
	_, err := db.Exec(query)
	return err
}
