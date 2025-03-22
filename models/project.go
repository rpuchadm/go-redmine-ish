package models

import (
	"database/sql"
	"log"
	"time"
)

// Project representa la estructura de un proyecto en la base de datos
type Project struct {
	ID          int           `json:"id"`
	ParentID    sql.NullInt64 `json:"parent_id"` // Usamos sql.NullInt64 para manejar valores NULL
	Name        string        `json:"name"`
	Identifier  string        `json:"identifier"`
	Description string        `json:"description"`
	CreatedOn   time.Time     `json:"created_on"`
	UpdatedOn   time.Time     `json:"updated_on"`
}

// CreateProject inserta un nuevo proyecto en la base de datos
func CreateProject(db *sql.DB, project *Project) error {
	query := `
	INSERT INTO projects (name, identifier, description, parent_id, created_on, updated_on)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	err := db.QueryRow(
		query,
		project.Name,
		project.Identifier,
		project.Description,
		project.ParentID,
		time.Now(),
		time.Now(),
	).Scan(&project.ID)

	if err != nil {
		log.Printf("Error al crear el proyecto: %v", err)
		return err
	}

	return nil
}

// GetProjectByID obtiene un proyecto por su ID
func GetProjectByID(db *sql.DB, id int) (*Project, error) {
	query := `
	SELECT id, name, identifier, description, parent_id, created_on, updated_on
	FROM projects
	WHERE id = $1`

	project := &Project{}
	var parentID sql.NullInt64

	err := db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Identifier,
		&project.Description,
		&parentID,
		&project.CreatedOn,
		&project.UpdatedOn,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Proyecto no encontrado
		}
		log.Printf("Error al obtener el proyecto: %v", err)
		return nil, err
	}

	project.ParentID = parentID
	return project, nil
}

// UpdateProject actualiza un proyecto existente en la base de datos
func UpdateProject(db *sql.DB, project *Project) error {
	query := `
	UPDATE projects
	SET name = $1, identifier = $2, description = $3, parent_id = $4, updated_on = $5
	WHERE id = $6`

	_, err := db.Exec(
		query,
		project.Name,
		project.Identifier,
		project.Description,
		project.ParentID,
		time.Now(),
		project.ID,
	)

	if err != nil {
		log.Printf("Error al actualizar el proyecto: %v", err)
		return err
	}

	return nil
}

// DeleteProject elimina un proyecto por su ID
func DeleteProject(db *sql.DB, id int) error {
	query := `DELETE FROM projects WHERE id = $1`

	_, err := db.Exec(query, id)
	if err != nil {
		log.Printf("Error al eliminar el proyecto: %v", err)
		return err
	}

	return nil
}
