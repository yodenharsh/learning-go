package models

import "database/sql"

type Exec struct {
	Id                    int            `json:"id,omitempty"`
	FirstName             string         `json:"firstName,omitempty"`
	LastName              string         `json:"lastName,omitempty"`
	Email                 string         `json:"email,omitempty"`
	Username              string         `json:"username,omitempty"`
	Password              string         `json:"password,omitempty"`
	InactiveStatus        bool           `json:"inactiveStatus,omitempty"`
	Role                  string         `json:"role,omitempty"`
	PasswordResetCode     sql.NullString `json:"passwordResetCode"`
	PasswordCodeExpiresAt sql.NullString `json:"passwordCodeExpiresAt"`
	PasswordChangedAt     sql.NullString `json:"passwordChangedAt"`
	UserCreatedAt         sql.NullString `json:"userCreatedAt"`
}
