package dbhelper

import (
	db "TODO/database"
	model "TODO/models"
	"database/sql"
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

func SuspendUser(userID string) (err error) {
	tx, err := db.Todo.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	query := `
		UPDATE users
		SET is_suspended = true
		WHERE id = $1
	`

	res, err := tx.Exec(query, userID)
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

	query = `
		UPDATE user_session
		SET archived_at = NOW()
		WHERE user_id = $1 AND archived_at IS NULL
	`

	_, err = tx.Exec(query, userID)
	return err
}

func UnsuspendUser(userID string) error {

	query := `
		UPDATE users
		SET is_suspended = false
		WHERE id = $1
	`
	_, err := db.Todo.Exec(query, userID)

	return err
}

func FetchAllUsers(limit, offset int) ([]model.User, error) {
	var users []model.User

	SQL := `SELECT id, name, email, created_at,role
        FROM users
        WHERE archived_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2;`

	err := db.Todo.Select(&users, SQL, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}
