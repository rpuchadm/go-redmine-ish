package models

import "database/sql"

/*
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,              -- Identificador único del comentario
    issue_id INT NOT NULL,              -- ID del ticket al que pertenece el comentario
    user_id INT NOT NULL,               -- ID del usuario que hizo el comentario
    content TEXT NOT NULL,              -- Contenido del comentario
    created_at TIMESTAMP DEFAULT NOW(), -- Fecha de creación del comentario
    updated_at TIMESTAMP DEFAULT NOW(), -- Fecha de última actualización del comentario
    FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
*/

type Comment struct {
	ID        int    `json:"id"`
	IssueID   int    `json:"issue_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateComment crea un nuevo comentario
func CreateComment(db *sql.DB, comment *Comment) (int, error) {
	query := `
	INSERT INTO comments (issue_id, user_id, content)
	VALUES ($1, $2, $3)
	RETURNING id`
	var id int
	err := db.QueryRow(query, comment.IssueID, comment.UserID, comment.Content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetCommentByID obtiene un comentario por su ID
func GetCommentByID(db *sql.DB, id int) (*Comment, error) {
	query := `
	SELECT id, issue_id, user_id, content, created_at, updated_at
	FROM comments
	WHERE id = $1`
	comment := &Comment{}
	err := db.QueryRow(query, id).Scan(&comment.ID, &comment.IssueID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentsByIssueID obtiene todos los comentarios de un ticket
func GetCommentsByIssueID(db *sql.DB, issueID int) ([]Comment, error) {
	query := `
	SELECT id, issue_id, user_id, content, created_at, updated_at
	FROM comments
	WHERE issue_id = $1`
	rows, err := db.Query(query, issueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.IssueID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

// UpdateComment actualiza un comentario
func UpdateComment(db *sql.DB, comment *Comment) error {
	query := `
	UPDATE comments
	SET content = $1, updated_at = NOW()
	WHERE id = $2`
	_, err := db.Exec(query, comment.Content, comment.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteComment elimina un comentario
func DeleteComment(db *sql.DB, id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// CountComments cuenta el número de comentarios
func CountComments(db *sql.DB) (int, error) {
	query := `SELECT COUNT(*) FROM comments`
	var count int
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// DropCommentsTable elimina la tabla de comentarios
func DropCommentsTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS comments`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func CreateCommentsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		issue_id INT NOT NULL,
		user_id INT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		FOREIGN KEY (issue_id) REFERENCES issues(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	)`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func SampleComments(db *sql.DB) error {
	query := `
	INSERT INTO comments (issue_id, user_id, content)
	VALUES
		(2, 2, 'Este es un comentario de prueba'),
		(3, 3, 'Este es otro comentario de prueba'),
		(2, 3, 'Este es un comentario de prueba para otro ticket')`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
