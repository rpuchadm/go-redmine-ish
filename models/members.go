package models

import "database/sql"

/*
CREATE TABLE IF NOT EXISTS members (
	id SERIAL PRIMARY KEY,
	user_id INT,
	project_id INT,
	role_id INT,
	created_at TIMESTAMP,
	updated_at TIMESTAMP
);
*/

type Member struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ProjectID int    `json:"project_id"`
	RoleID    int    `json:"role_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetAllMembers obtiene todos los miembros
func GetAllMembers(db *sql.DB) ([]Member, error) {
	query := `SELECT id, user_id, project_id, role_id, created_at, updated_at FROM members`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var member Member
		err := rows.Scan(&member.ID, &member.UserID, &member.ProjectID, &member.RoleID, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

// GetMemberByID obtiene un miembro por su ID
func GetMemberByID(db *sql.DB, id int) (*Member, error) {
	query := `SELECT id, user_id, project_id, role_id, created_at, updated_at FROM members WHERE id = $1`

	var member Member
	err := db.QueryRow(query, id).Scan(&member.ID, &member.UserID, &member.ProjectID, &member.RoleID, &member.CreatedAt, &member.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &member, nil
}

// GetMembersByProjectID obtiene todos los miembros de un proyecto
func GetMembersByProjectID(db *sql.DB, projectID int) ([]Member, error) {
	query := `SELECT id, user_id, project_id, role_id, created_at, updated_at FROM members WHERE project_id = $1`

	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var member Member
		err := rows.Scan(&member.ID, &member.UserID, &member.ProjectID, &member.RoleID, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

// GetMembersByUserID obtiene todos los miembros de un usuario
func GetMembersByUserID(db *sql.DB, userID int) ([]Member, error) {
	query := `SELECT id, user_id, project_id, role_id, created_at, updated_at FROM members WHERE user_id = $1`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var member Member
		err := rows.Scan(&member.ID, &member.UserID, &member.ProjectID, &member.RoleID, &member.CreatedAt, &member.UpdatedAt)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

// CreateMember crea un nuevo miembro
func CreateMember(db *sql.DB, member *Member) (int, error) {
	query := `INSERT INTO members (user_id, project_id, role_id) VALUES ($1, $2, $3) RETURNING id`

	var id int
	err := db.QueryRow(query, member.UserID, member.ProjectID, member.RoleID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateMember actualiza un miembro
func UpdateMember(db *sql.DB, member *Member) error {
	query := `UPDATE members SET user_id = $1, project_id = $2, role_id = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`

	_, err := db.Exec(query, member.UserID, member.ProjectID, member.RoleID, member.ID)
	return err
}

// DeleteMember elimina un miembro
func DeleteMember(db *sql.DB, id int) error {
	query := `DELETE FROM members WHERE id = $1`

	_, err := db.Exec(query, id)
	return err
}

// DeleteMembersByProjectID elimina todos los miembros de un proyecto
func DeleteMembersByProjectID(db *sql.DB, projectID int) error {
	query := `DELETE FROM members WHERE project_id = $1`

	_, err := db.Exec(query, projectID)
	return err
}

// DeleteMembersByUserID elimina todos los miembros de un usuario
func DeleteMembersByUserID(db *sql.DB, userID int) error {
	query := `DELETE FROM members WHERE user_id = $1`

	_, err := db.Exec(query, userID)
	return err
}

func CreateTableMembers(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS members (
		id SERIAL PRIMARY KEY,
		user_id INT,
		project_id INT,
		role_id INT,
		created_at TIMESTAMP,
		updated_at TIMESTAMP
	)`

	_, err := db.Exec(query)
	return err
}

func DropTableMembers(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS members`
	_, err := db.Exec(query)
	return err
}
