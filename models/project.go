package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Project representa la estructura de un proyecto en la base de datos
type Project struct {
	ID          int       `json:"id"`
	ParentID    *int      `json:"parent_id"` // Usamos *int para manejar valores NULL
	Name        string    `json:"name"`
	Identifier  string    `json:"identifier"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedOn   time.Time `json:"updated_on"`
}

// CreateProject inserta un nuevo proyecto en la base de datos
func CreateProject(db *sql.DB, project *Project) (int, error) {
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
		return 0, err
	}

	return project.ID, nil
}

// GetProjectByID obtiene un proyecto por su ID
func GetProjectByID(db *sql.DB, id int) (*Project, error) {
	query := `
	SELECT id, name, identifier, description, parent_id, created_on, updated_on
	FROM projects
	WHERE id = $1`

	project := &Project{}

	err := db.QueryRow(query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Identifier,
		&project.Description,
		&project.ParentID,
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

	return project, nil
}

func GetProjectsByUserID(db *sql.DB, userID int) ([]*Project, error) {
	query := `
	SELECT id, name, identifier, description, parent_id, created_on, updated_on
	FROM projects
	WHERE id IN (
		SELECT project_id
		FROM issues
		WHERE assigned_to_id = $1
	) or id IN (
	 	select project_id
		from categories
		where assigned_to_id = $1
	)
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Printf("Error al obtener los proyectos: %v", err)
		return nil, err
	}
	defer rows.Close()

	projects := []*Project{}
	for rows.Next() {
		project := &Project{}

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Identifier,
			&project.Description,
			&project.ParentID,
			&project.CreatedOn,
			&project.UpdatedOn,
		)
		if err != nil {
			log.Printf("Error al escanear el proyecto: %v", err)
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
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

// GetProjects obtiene todos los proyectos de la base de datos
func GetAllProjects(db *sql.DB) ([]*Project, error) {
	query := `
	SELECT id, name, identifier, description, parent_id, created_on, updated_on
	FROM projects
	ORDER BY id`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error al obtener los proyectos: %v", err)
		return nil, err
	}
	defer rows.Close()

	projects := []*Project{}
	for rows.Next() {
		project := &Project{}

		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Identifier,
			&project.Description,
			&project.ParentID,
			&project.CreatedOn,
			&project.UpdatedOn,
		)
		if err != nil {
			log.Printf("Error al escanear el proyecto: %v", err)
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func CountProjects(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM projects`

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error al contar los proyectos: %v", err)
		return 0, err
	}

	return count, nil
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

func CreateProjectsTable(db *sql.DB) error {
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

	_, err := db.Exec(createTableQuery)
	if err != nil {
		log.Printf("Error al crear la tabla projects: %v", err)
		return err
	}

	return nil
}

func DropProjectsTable(db *sql.DB) error {
	dropTableQuery := `DROP TABLE IF EXISTS projects;`

	_, err := db.Exec(dropTableQuery)
	if err != nil {
		log.Printf("Error al eliminar la tabla projects: %v", err)
		return err
	}

	return nil
}

func TestProjectsTable(db *sql.DB) error {
	// Crear un proyecto de prueba
	project0 := &Project{
		Name:        "Proyecto de ejemplo",
		Identifier:  "ejemplo",
		Description: "Este es un proyecto de ejemplo",
	}
	project_id, err := CreateProject(db, project0)
	if err != nil {
		fmt.Printf("Error al crear el proyecto: %v\n", err)
		return err
	}

	// cargar el proyecto creado
	project1, err := GetProjectByID(db, project_id)
	if err != nil {
		fmt.Printf("Error al obtener el proyecto: %v\n", err)
		return err
	}
	// comprobar que el proyecto se ha creado correctamente
	if project1.Name != project0.Name || project1.Identifier != project0.Identifier || project1.Description != project0.Description {
		fmt.Println("El proyecto no se ha creado correctamente")
		return fmt.Errorf("el proyecto no se ha creado correctamente")
	}
	// actualizar el proyecto creado
	project1.Name = "Proyecto de ejemplo actualizado"
	project1.Identifier = "ejemplo-actualizado"
	project1.Description = "Este es un proyecto de ejemplo actualizado"
	err = UpdateProject(db, project1)
	if err != nil {
		fmt.Printf("Error al actualizar el proyecto: %v\n", err)
		return err
	}
	// cargar el proyecto actualizado
	project2, err := GetProjectByID(db, project_id)
	if err != nil {
		fmt.Printf("Error al obtener el proyecto: %v\n", err)
		return err
	}
	// comprobar que el proyecto se ha actualizado correctamente
	if project2.Name != project1.Name || project2.Identifier != project1.Identifier || project2.Description != project1.Description {
		fmt.Println("El proyecto no se ha actualizado correctamente")
		return fmt.Errorf("el proyecto no se ha actualizado correctamente")
	}
	// eliminar el proyecto creado
	err = DeleteProject(db, project_id)
	if err != nil {
		fmt.Printf("Error al eliminar el proyecto: %v\n", err)
		return err
	}
	// comprobar que el proyecto se ha eliminado correctamente
	project3, err := GetProjectByID(db, project_id)
	if err != nil {
		fmt.Printf("El proyecto no se ha eliminado correctamente %v \n", err)
		return err
	}
	if project3 != nil {
		fmt.Println("El proyecto no se ha eliminado correctamente")
		return fmt.Errorf("el proyecto no se ha eliminado correctamente")
	}

	return nil
}

func SampleProjects(db *sql.DB) error {

	// insertamos 3 proyectos de ejemplo
	project4 := &Project{
		Name:        "Proyecto 1",
		Identifier:  "proyecto-1",
		Description: "Este es el proyecto 1",
	}
	_, err := CreateProject(db, project4)
	if err != nil {
		fmt.Printf("Error al crear el proyecto: %v\n", err)
		return err
	}

	project5 := &Project{
		Name:        "Proyecto 2",
		Identifier:  "proyecto-2",
		Description: "Este es el proyecto 2",
	}
	_, err = CreateProject(db, project5)
	if err != nil {
		fmt.Printf("Error al crear el proyecto: %v\n", err)
		return err
	}

	project6 := &Project{
		Name:        "Proyecto 3",
		Identifier:  "proyecto-3",
		Description: "Este es el proyecto 3",
	}
	_, err = CreateProject(db, project6)
	if err != nil {
		fmt.Printf("Error al crear el proyecto: %v\n", err)
		return err
	}

	return nil
}
