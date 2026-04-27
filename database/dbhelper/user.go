package dbhelper

import (
	db "TODO/database"
	models "TODO/models"
	"TODO/utils"
	"database/sql"
	"errors"
)

func IsUserExist(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)` //count() >0
	err := db.Todo.Get(&exists, query, email)
	return exists, err
}

func CreateUser(name, email, password string) (string, error) {
	query := `INSERT INTO users(name, email, password)
	VALUES ($1, TRIM(LOWER($2)), $3) RETURNING id;`

	var userID string
	err := db.Todo.Get(&userID, query, name, email, password)
	return userID, err
}

func CreateUserSession(userID string) (string, error) {
	query := `INSERT INTO user_session(user_id)
			VALUES ($1) RETURNING id;`
	var sessionID string
	err := db.Todo.Get(&sessionID, query, userID)
	if err != nil {
		return "", err
	}
	return sessionID, nil
}

func GetUserByEmail(email, password string) (models.User, error) {

	query := `
		SELECT id, password, role, is_suspended
		FROM users
		WHERE email = $1 AND archived_at IS NULL;
	`

	var user models.User

	err := db.Todo.Get(&user, query, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("no user exist")
		}
		return models.User{}, err
	}

	if err := utils.CheckPassword(user.Password, password); err != nil {
		return models.User{}, errors.New("invalid credentials")
	}

	return user, nil
}

func DeleteUserSession(sessionID string) error {
	query := `
	UPDATE user_session   
	SET archived_at = NOW()
	WHERE id = $1 AND archived_at IS NULL
`
	_, err := db.Todo.Exec(query, sessionID)
	return err
}

func GetTodoByID(userID, todoID string) (models.Todo, error) {
	query := `
		SELECT id, name, description, is_complete, expiring_at, created_at
		FROM todos
		WHERE id = $1 AND user_id = $2 AND archived_at IS NULL
	`

	var todo models.Todo

	err := db.Todo.Get(&todo, query, todoID, userID)
	return todo, err
}

func GetUserIDBySession(sessionID string) (string, error) {
	var userID string

	query := `
		SELECT user_id 
		FROM user_session 
		WHERE id = $1 AND archived_at IS NULL
	`

	err := db.Todo.Get(&userID, query, sessionID)
	return userID, err
}

func GetUserByID(userID string) (models.User, error) {

	var user models.User

	query := `
		SELECT id, name, email, role, is_suspended
		FROM users
		WHERE id = $1 AND archived_at IS NULL
	`

	err := db.Todo.Get(&user, query, userID)
	return user, err
}
