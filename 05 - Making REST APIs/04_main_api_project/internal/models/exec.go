package models

import "database/sql"

type Exec struct {
	Id                    int            `json:"id,omitempty" db:"id"`
	FirstName             string         `json:"firstName,omitempty" db:"first_name"`
	LastName              string         `json:"lastName,omitempty" db:"last_name"`
	Email                 string         `json:"email,omitempty" db:"email"`
	Username              string         `json:"username,omitempty" db:"username"`
	Password              string         `json:"password,omitempty" db:"password"`
	InactiveStatus        bool           `json:"inactiveStatus,omitempty" db:"inactive_status"`
	Role                  string         `json:"role,omitempty" db:"role"`
	PasswordResetCode     sql.NullString `json:"passwordResetCode" db:"password_reset_code"`
	PasswordCodeExpiresAt sql.NullString `json:"passwordCodeExpiresAt" db:"password_code_expires_at"`
	PasswordChangedAt     sql.NullString `json:"passwordChangedAt" db:"password_changed_at"`
	UserCreatedAt         sql.NullString `json:"userCreatedAt" db:"user_created_at"`
}

type UpdatePasswordRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}
