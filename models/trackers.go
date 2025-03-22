package models

import (
	"database/sql"
	"fmt"
)

/*
CREATE TABLE IF NOT EXISTS trackers (
    id SERIAL PRIMARY KEY,              -- Identificador único del tracker
    name VARCHAR(255) UNIQUE NOT NULL,  -- Nombre del tracker (por ejemplo, "Bug", "Feature")
    description TEXT                    -- Descripción del tracker
);
//Define los tipos de tickets o incidencias (por ejemplo, error, mejora, tarea).
*/

// Tracker representa un tipo de ticket o incidencia
type Tracker struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateTracker crea un nuevo tracker
func CreateTracker(db *sql.DB, tracker *Tracker) (int, error) {
	query := `INSERT INTO trackers (name, description) VALUES ($1, $2) RETURNING id`

	var id int
	err := db.QueryRow(query, tracker.Name, tracker.Description).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTrackerByID obtiene un tracker por su ID
func GetTrackerByID(db *sql.DB, id int) (*Tracker, error) {
	query := `SELECT id, name, description FROM trackers WHERE id = $1`

	tracker := &Tracker{}
	err := db.QueryRow(query, id).Scan(&tracker.ID, &tracker.Name, &tracker.Description)
	if err != nil {
		return nil, err
	}

	return tracker, nil
}

// GetTrackerByName obtiene un tracker por su nombre
func GetTrackerByName(db *sql.DB, name string) (*Tracker, error) {
	query := `SELECT id, name, description FROM trackers WHERE name = $1`

	tracker := &Tracker{}
	err := db.QueryRow(query, name).Scan(&tracker.ID, &tracker.Name, &tracker.Description)
	if err != nil {
		return nil, err
	}

	return tracker, nil
}

// GetAllTrackers obtiene todos los trackers
func GetAllTrackers(db *sql.DB) ([]*Tracker, error) {
	query := `SELECT id, name, description FROM trackers`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	trackers := []*Tracker{}
	for rows.Next() {
		tracker := &Tracker{}

		err := rows.Scan(
			&tracker.ID,
			&tracker.Name,
			&tracker.Description,
		)
		if err != nil {
			return nil, err
		}

		trackers = append(trackers, tracker)
	}

	return trackers, nil
}

// DropTrackersTable elimina la tabla de trackers
func DropTrackersTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS trackers`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// CountTrackers cuenta el número de trackers
func CountTrackers(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM trackers`

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateTracker actualiza un tracker
func UpdateTracker(db *sql.DB, tracker *Tracker) error {
	query := `UPDATE trackers SET name = $1, description = $2 WHERE id = $3`

	_, err := db.Exec(query, tracker.Name, tracker.Description, tracker.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTracker elimina un tracker
func DeleteTracker(db *sql.DB, id int) error {
	query := `DELETE FROM trackers WHERE id = $1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// CreateTrackersTable crea la tabla de trackers
func CreateTrackersTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS trackers (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		description TEXT
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

// SeedTrackers inserta datos de ejemplo en la tabla de trackers
func SeedTrackers(db *sql.DB) error {
	trackers := []*Tracker{
		{Name: "Bug", Description: "Error en el sistema"},
		{Name: "Feature", Description: "Nueva funcionalidad"},
		{Name: "Task", Description: "Tarea a realizar"},
	}

	for _, tracker := range trackers {
		_, err := CreateTracker(db, tracker)
		if err != nil {
			return err
		}
	}

	return nil
}

// TestTrackersTable verifica si la tabla de trackers existe
func TestTrackersTable(db *sql.DB) error {
	// Crear un tracker de prueba
	tracker0 := &Tracker{
		Name:        "Tracker de ejemplo",
		Description: "Este es un tracker de ejemplo",
	}

	tracker_id, err := CreateTracker(db, tracker0)
	if err != nil {
		return err
	}

	// cargar el tracker creado
	tracker1, err := GetTrackerByID(db, tracker_id)
	if err != nil {
		return err
	}

	// comprobar si el tracker creado es igual al tracker cargado
	if tracker0.Name != tracker1.Name || tracker0.Description != tracker1.Description {
		fmt.Printf("Error: los trackers no coinciden\n")
		return fmt.Errorf("los trackers no coinciden")
	}

	// modificar el tracker de prueba
	tracker1.Name = "Tracker modificado"
	tracker1.Description = "Este es un tracker modificado"

	err = UpdateTracker(db, tracker1)
	if err != nil {
		return err
	}

	// cargar el tracker modificado
	tracker2, err := GetTrackerByID(db, tracker_id)
	if err != nil {
		return err
	}

	// comprobar si el tracker modificado es igual al tracker cargado
	if tracker1.Name != tracker2.Name || tracker1.Description != tracker2.Description {
		fmt.Printf("Error: los trackers no coinciden\n")
		return fmt.Errorf("los trackers no coinciden")
	}

	// Eliminar el tracker de prueba
	err = DeleteTracker(db, tracker0.ID)
	if err != nil {
		return err
	}

	// comprobar si el tracker ha sido eliminado
	tracker3, err := GetTrackerByID(db, tracker_id)
	if err == nil {
		return err
	}

	if tracker3 != nil {
		fmt.Printf("Error: el tracker no se ha eliminado\n")
		return fmt.Errorf("el tracker no se ha eliminado")
	}

	return nil
}
