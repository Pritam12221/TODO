package models

import "time"

type UserRequest struct {
	Username string `json:"username" db:"username" binding:"required,min=3"`
	Password string `json:"password" db:"password" binding:"required,min=6"`
	Email    string `json:"email" db:"email" binding:"required,email"`
}

type User struct {
	Name        string     `json:"name" db:"name"`
	Password    string     `json:"password" db:"password"`
	ID          string     `json:"id" db:"id"`
	Email       string     `json:"email" db:"email"`
	Role        Role       `json:"role" db:"role"`
	IsSuspended bool       `json:"is_suspended" db:"is_suspended"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ArchivedAt  *time.Time `json:"archived_at" db:"archived_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type Todo struct {
	ID          string     `json:"id" db:"id"`
	Name        string     `json:"name" db:"name" binding:"required"`
	Description string     `json:"description" db:"description" binding:"required"`
	IsComplete  bool       `json:"is_complete" db:"is_complete"`
	ExpiringAt  *time.Time `json:"expiring_at" db:"expiring_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ArchivedAt  *time.Time `json:"archived_at" db:"archived_at"`
}

type UpdateTodoRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	IsComplete  *bool      `json:"is_complete"`
	ExpiringAt  *time.Time `json:"expiring_at" `
}

type UserExist struct {
	ID       string `db:"id"`
	Password string `db:"password"`
}

type AuthContext struct {
	UserID    string
	SessionID string
}

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleEmployee Role = "employee"
)

// type User struct {
// 	ID          string `db:"id"`
// 	Email       string `db:"email"`
// 	Role        Role   `db:"role"`
// 	IsSuspended bool   `db:"is_suspended"`
// }
