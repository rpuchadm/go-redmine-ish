package models

import (
	"database/sql"
	"fmt"
)

// User representa un usuario del sistema
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateUser crea un nuevo usuario
func CreateUser(db *sql.DB, user *User) (int, error) {
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`

	var id int
	err := db.QueryRow(query, user.Username, user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUserByID obtiene un usuario por su ID
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`

	user := &User{}
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByUsername obtiene un usuario por su nombre de usuario
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1`

	user := &User{}
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByEmail obtiene un usuario por su correo electrónico
func GetUserByEmail(db *sql.DB, email string) (*User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`

	user := &User{}
	err := db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser actualiza un usuario
func UpdateUser(db *sql.DB, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, updated_at = NOW() WHERE id = $4`

	_, err := db.Exec(query, user.Username, user.Email, user.PasswordHash, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser elimina un usuario
func DeleteUser(db *sql.DB, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// CountUsers cuenta el número de usuarios
func CountUsers(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetUsers obtiene todos los usuarios
func GetAllUsers(db *sql.DB) ([]*User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users`

	rows, err := db.Query(query)
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

// DropUsersTable elimina la tabla de usuarios
func DropUsersTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS users`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// CreateUsersTable crea la tabla de usuarios
func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW()
	);`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// TestUsersTable verifica si la tabla de usuarios existe
func TestUsersTable(db *sql.DB) error {

	// crea un user de prueba
	user := &User{
		Username:     "test",
		Email:        "test@mydomain.com",
		PasswordHash: "123456",
	}

	user_id, err := CreateUser(db, user)
	if err != nil {
		fmt.Printf("Error al crear el usuario de prueba: %v\n", err)
		return err
	}

	// cargar el user creado
	user1, err := GetUserByID(db, user_id)
	if err != nil {
		fmt.Printf("Error al obtener el usuario de prueba: %v\n", err)
		return err
	}

	// comprobar que el user se ha creado correctamente
	if user1.Username != user.Username || user1.Email != user.Email || user1.PasswordHash != user.PasswordHash {
		fmt.Printf("Error: los datos del usuario de prueba no coinciden %v %v\n", user1, user)
		return fmt.Errorf("error: los datos del usuario de prueba no coinciden")
	}

	// actualizar el user creado
	user1.Username = "test2"
	user1.Email = "test2@mydomain.com"
	user1.PasswordHash = "654321"

	err = UpdateUser(db, user1)
	if err != nil {
		fmt.Printf("Error al actualizar el usuario de prueba: %v\n", err)
		return err
	}

	// cargar el user actualizado
	user2, err := GetUserByID(db, user_id)
	if err != nil {
		fmt.Printf("Error al obtener el usuario de prueba actualizado: %v\n", err)
		return err
	}

	// comprobar que el user se ha actualizado correctamente
	if user2.Username != user1.Username || user2.Email != user1.Email || user2.PasswordHash != user1.PasswordHash {
		fmt.Printf("Error: los datos del usuario de prueba actualizado no coinciden\n")
		return fmt.Errorf("error: los datos del usuario de prueba actualizado no coinciden %v %v", user2, user1)
	}

	// eliminar el user creado
	err = DeleteUser(db, user_id)
	if err != nil {
		fmt.Printf("Error al eliminar el usuario de prueba: %v\n", err)
		return err
	}

	// comprobar que el user se ha eliminado correctamente
	user3, _ := GetUserByID(db, user_id)

	if user3 != nil {
		fmt.Printf("Error: el usuario de prueba no se ha eliminado correctamente\n")
		return fmt.Errorf("error: el usuario de prueba no se ha eliminado correctamente")
	}

	return nil
}

// SampleUsers crea usuarios de ejemplo
func SampleUsers(db *sql.DB) error {
	users := []*User{
		{
			Username:     "admin1",
			Email:        "admin1@mydomain.com",
			PasswordHash: "admin1",
		},
		{
			Username:     "user1",
			Email:        "user1@mydomain.com",
			PasswordHash: "user1",
		},
		{
			Username:     "user2",
			Email:        "user2@mydomain.com",
			PasswordHash: "user2",
		},
	}

	for _, user := range users {
		_, err := CreateUser(db, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetUsersByIssueID(db *sql.DB, issue_id int) ([]*User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users
	WHERE id IN (
		SELECT user_id
		FROM comments
		WHERE issue_id = $1
	) or id IN (
	 	select assigned_to_id
		from issues
		where id = $1
	)
	`

	rows, err := db.Query(query, issue_id)
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

func GetUsersByProjectID(db *sql.DB, project_id int) ([]*User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at, updated_at
	FROM users
	WHERE id IN (
		select assigned_to_id
		from categories
		where project_id = $1
	)
	`

	rows, err := db.Query(query, project_id)
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
