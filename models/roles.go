package models

import (
	"database/sql"
	"fmt"
)

// Role representa un rol que puede tener un usuario
type Role struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateRole crea un nuevo rol
func CreateRole(db *sql.DB, role *Role) (int, error) {
	query := `INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id`

	var id int
	err := db.QueryRow(query, role.Name, role.Description).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetRoleByID obtiene un rol por su ID
func GetRoleByID(db *sql.DB, id int) (*Role, error) {
	query := `SELECT id, name, description FROM roles WHERE id = $1`

	role := &Role{}
	err := db.QueryRow(query, id).Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetRoleByName obtiene un rol por su nombre
func GetRoleByName(db *sql.DB, name string) (*Role, error) {
	query := `SELECT id, name, description FROM roles WHERE name = $1`

	role := &Role{}
	err := db.QueryRow(query, name).Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// GetAllRoles obtiene todos los roles
func GetAllRoles(db *sql.DB) ([]Role, error) {
	query := `SELECT id, name, description FROM roles`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []Role{}
	for rows.Next() {
		role := Role{}

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

// UpdateRole actualiza un rol
func UpdateRole(db *sql.DB, role *Role) error {
	query := `UPDATE roles SET name = $1, description = $2 WHERE id = $3`

	_, err := db.Exec(query, role.Name, role.Description, role.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRole elimina un rol
func DeleteRole(db *sql.DB, id int) error {
	query := `DELETE FROM roles WHERE id = $1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// CountRoles cuenta el número de roles
func CountRoles(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM roles`

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func DropRolesTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS roles`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func CreateRolesTable(db *sql.DB) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS roles (
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

func SeedRolesTable(db *sql.DB) error {
	seedQuery := `
	INSERT INTO roles (name, description) VALUES
	('Admin', 'Administrador del sistema'),
	('Developer', 'Desarrollador de software'),
	('Reporter', 'Reportero de problemas');
	`

	_, err := db.Exec(seedQuery)
	if err != nil {
		return err
	}

	return nil
}

func TestRolesTable(db *sql.DB) error {

	// crear rol de prueba
	role := &Role{
		Name:        "Tester",
		Description: "Probador de software",
	}

	role_id, err := CreateRole(db, role)
	if err != nil {
		return err
	}

	// obtener rol por ID
	role1, err := GetRoleByID(db, role_id)
	if err != nil {
		return err
	}

	// comprobar si el rol es correcto
	if role1.Name != role.Name || role1.Description != role.Description {
		return fmt.Errorf("error: los datos del rol de prueba no coinciden %v %v", role1, role)
	}

	// actualizar rol
	role1.Name = "Tester 2"
	role1.Description = "Probador de software 2"
	err = UpdateRole(db, role1)
	if err != nil {
		return err
	}

	// obtener rol actualizado
	role2, err := GetRoleByID(db, role_id)
	if err != nil {
		return err
	}

	// comprobar si el rol se ha actualizado correctamente
	if role2.Name != role1.Name || role2.Description != role1.Description {
		return fmt.Errorf("error: los datos del rol de prueba actualizado no coinciden")
	}

	// eliminar rol de prueba
	err = DeleteRole(db, role_id)
	if err != nil {
		return err
	}

	// comprobar si el rol ha sido eliminado
	role3, _ := GetRoleByID(db, role_id)
	if role3 != nil {
		return fmt.Errorf("error: el rol no se ha eliminado")
	}

	return nil
}

// GetRolesByUserID obtiene los roles de un usuario
func GetRolesByUserID(db *sql.DB, userID int) ([]Role, error) {
	query := `
	SELECT id, name, description
	FROM roles r
	WHERE r.id IN (
		SELECT role_id
		FROM user_roles
		WHERE user_id = $1
	)`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []Role{}
	for rows.Next() {
		role := Role{}

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

/*
// AddRoleToUser agrega un rol a un usuario
func AddRoleToUser(db *sql.DB, userID, roleID int) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`

	_, err := db.Exec(query, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveRoleFromUser elimina un rol de un usuario
func RemoveRoleFromUser(db *sql.DB, userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	_, err := db.Exec(query, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAllRolesFromUser elimina todos los roles de un usuario
func RemoveAllRolesFromUser(db *sql.DB, userID int) error {
	query := `DELETE FROM user_roles WHERE user_id = $1`

	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}

// GetUsersByRoleID obtiene los usuarios con un rol
func GetUsersByRoleID(db *sql.DB, roleID int) ([]*User, error) {
	query := `
	SELECT u.id, u.username, u.email, u.password_hash, u.created_at, u.updated_at
	FROM users u
	JOIN user_roles ur ON u.id = ur.user_id
	WHERE ur.role_id = $1`

	rows, err := db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := &User{}

		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetRolesByProjectID obtiene los roles de un proyecto
func GetRolesByProjectID(db *sql.DB, projectID int) ([]*Role, error) {
	query := `
	SELECT r.id, r.name, r.description
	FROM roles r
	JOIN project_roles pr ON r.id = pr.role_id
	WHERE pr.project_id = $1`

	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*Role{}
	for rows.Next() {
		role := &Role{}

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

// AddRoleToProject agrega un rol a un proyecto
func AddRoleToProject(db *sql.DB, projectID, roleID int) error {
	query := `INSERT INTO project_roles (project_id, role_id) VALUES ($1, $2)`

	_, err := db.Exec(query, projectID, roleID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveRoleFromProject elimina un rol de un proyecto
func RemoveRoleFromProject(db *sql.DB, projectID, roleID int) error {
	query := `DELETE FROM project_roles WHERE project_id = $1 AND role_id = $2`

	_, err := db.Exec(query, projectID, roleID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAllRolesFromProject elimina todos los roles de un proyecto
func RemoveAllRolesFromProject(db *sql.DB, projectID int) error {
	query := `DELETE FROM project_roles WHERE project_id = $1`

	_, err := db.Exec(query, projectID)
	if err != nil {
		return err
	}

	return nil
}

// GetProjectsByRoleID obtiene los proyectos con un rol
func GetProjectsByRoleID(db *sql.DB, roleID int) ([]*Project, error) {
	query := `
	SELECT p.id, p.name, p.identifier, p.description, p.parent_id, p.created_on, p.updated_on
	FROM projects p
	JOIN project_roles pr ON p.id = pr.project_id
	WHERE pr.role_id = $1`

	rows, err := db.Query(query, roleID)
	if err != nil {
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
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// GetRolesByIssueID obtiene los roles de un problema
func GetRolesByIssueID(db *sql.DB, issueID int) ([]*Role, error) {
	query := `
	SELECT r.id, r.name, r.description
	FROM roles r
	JOIN issue_roles ir ON r.id = ir.role_id
	WHERE ir.issue_id = $1`

	rows, err := db.Query(query, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*Role{}
	for rows.Next() {
		role := &Role{}

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
		)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

// AddRoleToIssue agrega un rol a un problema
func AddRoleToIssue(db *sql.DB, issueID, roleID int) error {
	query := `INSERT INTO issue_roles (issue_id, role_id) VALUES ($1, $2)`

	_, err := db.Exec(query, issueID, roleID)
	if err != nil {
		return err
	}

	return nil
}
*/
