package dbhelper

import (
	db "TODO/database"
	model "TODO/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

// create a todo
func CreateTodo(userID, name, description string, expiringAt *time.Time) (model.Todo, error) {
	query := `
		INSERT INTO todos (user_id, name, description, expiring_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, is_complete, expiring_at, created_at
	`
	var todo model.Todo
	err := db.Todo.Get(&todo, query, userID, name, description, expiringAt)
	return todo, err
}

// updates the todo
func UpdateTodoRequest(todoID, userID string, req model.UpdateTodoRequest) error {
	query := `
  UPDATE todos
  SET
   name        = COALESCE($1, name),
   description = COALESCE($2, description),
   is_complete    = COALESCE($3, is_complete),
   expiring_at = COALESCE($4, expiring_at)
  WHERE id = $5 AND user_id = $6
 `
	res, err := db.Todo.Exec(
		query,
		req.Name,
		req.Description,
		req.IsComplete,
		req.ExpiringAt,
		todoID,
		userID,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func DeleteTodo(todoID, userID string) error {
	query := `
		UPDATE todos
		SET archived_at = NOW()
		WHERE id = $1 
		  AND user_id = $2 
		  AND archived_at IS NULL
	`
	result, err := db.Todo.Exec(query, todoID, userID)

	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func GetTodoById(userID, todoID string) (model.Todo, error) {
	query := `
		SELECT id, name, description, is_complete, expiring_at, created_at
		FROM todos
		WHERE id = $1 AND user_id = $2 AND archived_at IS NULL
	`

	var todo model.Todo
	log.Println("todoid", todoID)

	err := db.Todo.Get(&todo, query, todoID, userID)
	return todo, err
}

func GetTodosByStatus(userID string, status string, search string, limit int, offset int) ([]model.Todo, error) {

	query := `
		SELECT id, name, description, is_complete, expiring_at, created_at, archived_at
		FROM todos
		WHERE user_id = $1
		  AND archived_at IS NULL
	`

	args := []any{userID}
	i := 2

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

// func GetTodosByStatus(userID string, status model.TodoStatus, search string, limit int, offset int) ([]model.Todo, error) {

// 	// Optional validation
// 	if status != "" &&
// 		status != model.StatusCompleted &&
// 		status != model.StatusPending &&
// 		status != model.StatusIncomplete {
// 		return nil, errors.New("invalid status")
// 	}

// 	query := `
// 		SELECT id, name, description, is_complete, expiring_at, created_at, archived_at
// 		FROM todos
// 		WHERE user_id = $1
// 		  AND archived_at IS NULL
// 		  AND (
// 			$2 = '' OR
// 			($2 = 'completed' AND is_complete = true) OR
// 			($2 = 'pending' AND is_complete = false AND (expiring_at IS NULL OR expiring_at > NOW())) OR
// 			($2 = 'incomplete' AND is_complete = false AND (expiring_at IS NULL OR expiring_at < NOW()))
// 		  )
// 	`

// 	args := []any{userID, string(status)}
// 	i := 3

// 	if search != "" {
// 		query += fmt.Sprintf(`
// 			AND (name ILIKE $%d OR description ILIKE $%d)
// 		`, i, i)
// 		args = append(args, "%"+search+"%")
// 		i++
// 	}

// 	query += fmt.Sprintf(`
// 		ORDER BY created_at DESC
// 		LIMIT $%d OFFSET $%d
// 	`, i, i+1)

// 	args = append(args, limit, offset)

// 	var todos []model.Todo
// 	err := db.Todo.Select(&todos, query, args...)
// 	return todos, err
// }
