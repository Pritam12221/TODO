package dbhelper

import (
	db "TODO/database"
	model "TODO/models"
	"errors"
	"fmt"
)

func GetAllTodos(status string, search string, limit int, offset int) ([]model.Todo, error) {

	query := `
		SELECT id, name, description, is_complete, expiring_at, created_at, archived_at
		FROM todos
		WHERE archived_at IS NULL
	`

	args := []any{}
	i := 1

	if search != "" {
		query += fmt.Sprintf(`
			AND (name ILIKE $%d OR description ILIKE $%d)
		`, i, i)
		args = append(args, "%"+search+"%")
		i++
	}

	switch status {
	case "incomplete":
		query += `
			AND is_complete = false 
			AND (expiring_at IS NULL OR expiring_at < NOW())
		`

	case "pending":
		query += `
			AND is_complete = false 
			AND (expiring_at IS NULL OR expiring_at > NOW())
		`

	case "completed":
		query += `
			AND is_complete = true
		`
		//get all todo
	case "":

	default:
		return nil, errors.New("invalid status")
	}

	query += fmt.Sprintf(`
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, i, i+1)

	args = append(args, limit, offset)

	var todos = make([]model.Todo, 0)
	err := db.Todo.Select(&todos, query, args...)
	return todos, err
}

func SuspendUser(userID string) error {
	query := `
		UPDATE users
		SET is_suspended = true
		WHERE id = $1
	`
	_, err := db.Todo.Exec(query, userID)
	return err
}
