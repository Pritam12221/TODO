package dbhelper

import (
	db "TODO/database"
	model "TODO/models"
	"errors"
	"log"
	"time"
)

func CreateTodo(userID, name, description string, expiringAt *time.Time) (model.Todo, error) {
	query := `
		INSERT INTO todos (user_id, name, description, expiring_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, complete, expiring_at, created_at
	`
	var todo model.Todo
	err := db.Todo.QueryRow(query, userID, name, description, expiringAt).Scan(
		&todo.ID, &todo.Name, &todo.Description, &todo.Complete, &todo.ExpiringAt, &todo.CreatedAt,
	)
	return todo, err
}

func UpdateTodoRequest(todoID, userID string, req model.UpdateTodoRequest) error {
	query := `
		UPDATE todos
		SET
			name        = $1,
			description = $2,
			complete    = $3,
			expiring_at = $4
		WHERE id = $5 
		  AND user_id = $6 
		  AND archived_at IS NULL
	`

	result, err := db.Todo.Exec(
		query,
		req.Name,
		req.Description,
		req.Complete,
		req.ExpiringAt,
		todoID,
		userID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("todo not found or already deleted")
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
 rowsAffected, err := result.RowsAffected()
 if err != nil {
  return err
 }
 if rowsAffected == 0 {
  return errors.New("todo not found")
 }
 return nil
}

func GetTodoById(userID, todoID string) (model.Todo, error) {
	query := `
		SELECT id, name, description, complete, expiring_at, created_at
		FROM todos
		WHERE id = $1 AND user_id = $2 AND archived_at IS NULL
	`

	var todo model.Todo
	log.Println("todoid", todoID)
	
	err := db.Todo.Get(&todo, query, todoID, userID)
	return todo, err
}



func GetTodosByStatus(userID, status string) ([]model.Todo, error) {
	baseQuery := `
		SELECT id, name, description, complete, expiring_at, created_at, archived_at
		FROM todos
		WHERE user_id = $1
		  AND archived_at IS NULL
	`
log.Println("kam kargyi")
	var query string

	switch status {
	case "incomplete":
		query = baseQuery + " AND complete = false AND (expiring_at IS NULL OR expiring_at < NOW()) ORDER BY created_at DESC"

	case "pending":
		query = baseQuery + " AND complete = false AND (expiring_at IS NULL OR expiring_at > NOW()) ORDER BY created_at DESC"

	case "completed":
		query = baseQuery + " AND complete = true ORDER BY created_at DESC"

	case "":
		query = baseQuery + " ORDER BY created_at DESC"

	default:
		return nil, errors.New("invalid status")
	}

	var todos []model.Todo
	err := db.Todo.Select(&todos, query, userID)
	return todos, err
}