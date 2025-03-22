package models

import "database/sql"

// UserRole representa la relación entre un usuario y un rol
type UserRole struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
}

// CreateUserRoles crea una nueva relación entre un usuario y un rol
func CreateUserRoles(db *sql.DB, userRole *UserRole) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`

	_, err := db.Exec(query, userRole.UserID, userRole.RoleID)
	if err != nil {
		return err
	}

	return nil
}

// GetUserRolesByUserID obtiene los roles de un usuario por su ID
func GetUserRolesByUserID(db *sql.DB, userID int) ([]*Role, error) {
	query := `SELECT r.id, r.name, r.description FROM roles r
	JOIN user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*Role{}
	for rows.Next() {
		role := &Role{}
		err := rows.Scan(&role.ID, &role.Name, &role.Description)
		if err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	return roles, nil
}

// DeleteUserRoles elimina todas las relaciones de un usuario
func DeleteUserRoles(db *sql.DB, userID int) error {
	query := `DELETE FROM user_roles WHERE user_id = $1`

	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserRole elimina una relación entre un usuario y un rol
func DeleteUserRole(db *sql.DB, userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	_, err := db.Exec(query, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRoleUsers elimina todas las relaciones de un rol
func DeleteRoleUsers(db *sql.DB, roleID int) error {
	query := `DELETE FROM user_roles WHERE role_id = $1`

	_, err := db.Exec(query, roleID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRoleUser elimina una relación entre un rol y un usuario
func DeleteRoleUser(db *sql.DB, roleID, userID int) error {
	query := `DELETE FROM user_roles WHERE role_id = $1 AND user_id = $2`

	_, err := db.Exec(query, roleID, userID)
	if err != nil {
		return err
	}

	return nil
}

// CountUserRoles cuenta el número de roles de un usuario
func CountUserRoles(db *sql.DB, userID int) (int, error) {
	query := `SELECT COUNT(*) FROM user_roles WHERE user_id = $1`

	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
