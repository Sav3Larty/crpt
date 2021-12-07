package dbs

import (
	"database/sql"
	"forum/models"
)

// CreateCategory by using name.
func CreateCategory(name string) error {
	_, err := conn.Exec(
		"INSERT INTO category (name) VALUES(?)",
		name,
	)
	return err
}

// GetCategories with names of category.
func GetCategories() ([]*models.Category, error) {
	cats := make([]*models.Category, 0)

	rows, err := conn.Query("SELECT id, name FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		cat := &models.Category{}

		err = rows.Scan(&cat.ID, &cat.Name)

		if err != nil {
			return nil, err
		}
		cats = append(cats, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cats, nil
}

// ExistsCategory return if category exists.
func ExistsCategory(id int) (bool, error) {
	var exists bool
	err := conn.QueryRow("SELECT EXISTS (SELECT id FROM category WHERE id =?)", id).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return exists, err
	}
	return exists, nil
}
