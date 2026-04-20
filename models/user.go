package models

import "time"

type UserRequest struct {
	Username string `json:"username" db:"username" binding:"required,min=3"`
	Password string `json:"password" db:"password" binding:"required,min=6"`
	Email    string `json:"email" db:"email" binding:"required,email"`
}

type User struct {
	Name       string     `json:"name" db:"name"`
	Password   string     `json:"password" db:"password"`
	ID         string     `json:"id" db:"id"`
	Email      string     `json:"email" db:"email"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ArchivedAt *time.Time `json:"archived_at" db:"archived_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type Todo struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" binding:"required"`
	Description string     `json:"description" db:"description" binding:"required"`
	Complete    bool       `json:"complete" db:"complete"`
	ExpiringAt  *time.Time `json:"expiring_at" db:"expiring_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ArchivedAt  *time.Time `json:"archived_at" db:"archived_at"`
}

type UpdateTodoRequest struct {
 Name        *string    `json:"name"`
 Description *string    `json:"description"`
 Complete    *bool      `json:"complete"`
 ExpiringAt  *time.Time `json:"expiring_at" `
}

type UserExist struct {
	ID       string `db:"id"`
	Password string `db:"password"`
}
