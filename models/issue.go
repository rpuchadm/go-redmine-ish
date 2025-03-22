package models

import (
	"database/sql"
	"fmt"
)

// Issue representa un ticket o incidencia
type Issue struct {
	ID           int    `json:"id"`
	Subject      string `json:"subject"`
	Description  string `json:"description"`
	TrackerID    int    `json:"tracker_id"`
	ProjectID    int    `json:"project_id"`
	AssignedToID int    `json:"assigned_to_id"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateIssue crea un nuevo ticket
func CreateIssue(db *sql.DB, issue *Issue) (int, error) {
	query := `INSERT INTO issues (subject, description, tracker_id, project_id, assigned_to_id, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id int
	err := db.QueryRow(query, issue.Subject, issue.Description, issue.TrackerID, issue.ProjectID, issue.AssignedToID, issue.Status).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetIssueByID obtiene un ticket por su ID
func GetIssueByID(db *sql.DB, id int) (*Issue, error) {
	query := `SELECT id, subject, description, tracker_id, project_id, assigned_to_id, status, created_at, updated_at FROM issues WHERE id = $1`

	issue := &Issue{}
	err := db.QueryRow(query, id).Scan(&issue.ID, &issue.Subject, &issue.Description, &issue.TrackerID, &issue.ProjectID, &issue.AssignedToID, &issue.Status, &issue.CreatedAt, &issue.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return issue, nil
}

// GetIssuesByProjectID obtiene todos los tickets de un proyecto
func GetIssuesByProjectID(db *sql.DB, projectID int) ([]Issue, error) {
	query := `SELECT id, subject, description, tracker_id, project_id, assigned_to_id, status, created_at, updated_at FROM issues WHERE project_id = $1`

	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []Issue
	for rows.Next() {
		var issue Issue
		err := rows.Scan(
			&issue.ID,
			&issue.Subject,
			&issue.Description,
			&issue.TrackerID,
			&issue.ProjectID,
			&issue.AssignedToID,
			&issue.Status,
			&issue.CreatedAt,
			&issue.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// UpdateIssue actualiza un ticket existente en la base de datos
func UpdateIssue(db *sql.DB, issue *Issue) error {
	query := `UPDATE issues SET subject = $1, description = $2, tracker_id = $3, project_id = $4, assigned_to_id = $5, status = $6, updated_at = NOW() WHERE id = $7`

	_, err := db.Exec(query, issue.Subject, issue.Description, issue.TrackerID, issue.ProjectID, issue.AssignedToID, issue.Status, issue.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteIssue elimina un ticket
func DeleteIssue(db *sql.DB, id int) error {
	query := `DELETE FROM issues WHERE id = $1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// CountIssues cuenta el número de tickets
func CountIssues(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM issues`

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetAllIssues obtiene todos los tickets
func GetAllIssues(db *sql.DB) ([]Issue, error) {
	query := `SELECT id, subject, description, tracker_id, project_id, assigned_to_id, status, created_at, updated_at FROM issues`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []Issue
	for rows.Next() {
		var issue Issue
		err := rows.Scan(
			&issue.ID,
			&issue.Subject,
			&issue.Description,
			&issue.TrackerID,
			&issue.ProjectID,
			&issue.AssignedToID,
			&issue.Status,
			&issue.CreatedAt,
			&issue.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// DropIssuesTable elimina la tabla de tickets
func DropIssuesTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS issues`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// CreateIssuesTable crea la tabla de tickets
func CreateIssuesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS issues (
		id SERIAL PRIMARY KEY,
		subject VARCHAR(255) NOT NULL,
		description TEXT,
		tracker_id INT NOT NULL,
		project_id INT,
		assigned_to_id INT,
		status VARCHAR(50) DEFAULT 'Open',
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		FOREIGN KEY (tracker_id) REFERENCES trackers(id) ON DELETE SET NULL,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL,
		FOREIGN KEY (assigned_to_id) REFERENCES users(id) ON DELETE SET NULL
	)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// TestIssuesTable verifica si la tabla de tickets existe
func TestIssuesTable(db *sql.DB) error {
	// crea un ticket de prueba
	issue := &Issue{
		Subject:      "Test Issue",
		Description:  "This is a test issue",
		TrackerID:    2,
		ProjectID:    2,
		AssignedToID: 2,
		Status:       "Open",
	}

	issue_id, err := CreateIssue(db, issue)
	if err != nil {
		return err
	}

	// comprueba si el ticket se ha creado correctamente
	issue1, err := GetIssueByID(db, issue_id)
	if err != nil {
		return err
	}

	// comprueba si el ticket se ha creado correctamente
	if issue1.Subject != issue.Subject || issue1.Description != issue.Description {
		return fmt.Errorf("error: el ticket no se ha creado correctamente")
	}

	// actualiza el ticket de prueba
	issue1.Subject = "Test Issue Updated"
	err = UpdateIssue(db, issue1)
	if err != nil {
		return err
	}

	// comprueba si el ticket se ha actualizado correctamente
	issue2, err := GetIssueByID(db, issue_id)
	if err != nil {
		return err
	}

	// comprueba si el ticket se ha actualizado correctamente
	if issue2.Subject != issue1.Subject || issue2.Description != issue1.Description {
		return fmt.Errorf("error: el ticket no se ha actualizado correctamente")
	}

	// elimina el ticket de prueba
	err = DeleteIssue(db, issue_id)
	if err != nil {
		return err
	}

	return nil
}

// SeedIssues inserta datos de ejemplo en la tabla de tickets
func SampleIssues(db *sql.DB) error {

	issues := []*Issue{
		{Subject: "Issue 1", Description: "This is issue 1", TrackerID: 1, ProjectID: 2, AssignedToID: 2, Status: "Open"},
		{Subject: "Issue 2", Description: "This is issue 2", TrackerID: 2, ProjectID: 2, AssignedToID: 2, Status: "Open"},
		{Subject: "Issue 3", Description: "This is issue 3", TrackerID: 3, ProjectID: 2, AssignedToID: 2, Status: "Open"},
	}

	for _, issue := range issues {
		_, err := CreateIssue(db, issue)
		if err != nil {
			return err
		}
	}

	return nil
}
