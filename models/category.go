package models

import "database/sql"

type Category struct {
	ID           int    `json:"id"`
	ProjectID    int    `json:"project_id"`
	Name         string `json:"name"`
	AssignedToID *int   `json:"assigned_to_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateCategory crea una nueva categoría
func CreateCategory(db *sql.DB, category *Category) (int, error) {
	query := `
	INSERT INTO categories (project_id, name, assigned_to_id)
	VALUES ($1, $2, $3)
	RETURNING id`
	var id int
	err := db.QueryRow(query, category.ProjectID, category.Name, category.AssignedToID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetCategoryByID obtiene una categoría por su ID
func GetCategoryByID(db *sql.DB, id int) (*Category, error) {
	query := `
	SELECT id, project_id, name, assigned_to_id, created_at, updated_at
	FROM categories
	WHERE id = $1`
	category := &Category{}
	err := db.QueryRow(query, id).Scan(&category.ID, &category.ProjectID, &category.Name, &category.AssignedToID, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetCategoriesByProjectID obtiene todas las categorías de un proyecto
func GetCategoriesByProjectID(db *sql.DB, projectID int) ([]Category, error) {
	query := `
	SELECT id, project_id, name, assigned_to_id, created_at, updated_at
	FROM categories
	WHERE project_id = $1`
	rows, err := db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.ProjectID, &category.Name, &category.AssignedToID, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func GetCategoriesByUserID(db *sql.DB, userID int) ([]Category, error) {
	query := `
	SELECT id, project_id, name, assigned_to_id, created_at, updated_at
	FROM categories
	WHERE assigned_to_id = $1
	or id in (
		select category_id from issues
		where assigned_to_id = $1
	)
	`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var category Category
		if err := rows.Scan(&category.ID, &category.ProjectID, &category.Name, &category.AssignedToID, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// UpdateCategory actualiza una categoría
func UpdateCategory(db *sql.DB, category *Category) error {
	query := `
	UPDATE categories
	SET name = $1, assigned_to_id = $2, updated_at = NOW()
	WHERE id = $3`
	_, err := db.Exec(query, category.Name, category.AssignedToID, category.ID)
	return err
}

// DeleteCategory elimina una categoría
func DeleteCategory(db *sql.DB, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}

func CreateCategoriesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		project_id INT NOT NULL,
		name VARCHAR(255) NOT NULL,
		assigned_to_id INT,
		created_at TIMESTAMP DEFAULT NOW(),
		updated_at TIMESTAMP DEFAULT NOW(),
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
		FOREIGN KEY (assigned_to_id) REFERENCES users(id) ON DELETE SET NULL
	)`
	_, err := db.Exec(query)
	return err
}

func DropCategoriesTable(db *sql.DB) error {
	query := `DROP TABLE IF EXISTS categories`
	_, err := db.Exec(query)
	return err
}

func SampleCategories(db *sql.DB) error {
	query := `
	INSERT INTO categories (project_id, name, assigned_to_id)
	VALUES
	(4, 'General', NULL),
	(4, 'Desarrollo', 3),
	(4, 'Diseño', 4),
	(2, 'General', NULL),
	(2, 'Desarrollo', 3),
	(2, 'Diseño', 4),
	(3, 'General', NULL),
	(3, 'Desarrollo', 3),
	(3, 'Diseño', 4)`
	_, err := db.Exec(query)
	return err
}
