package models

type Student struct {
	Id        int    `json:"id,omitempty" db:"id,omitempty"`
	FirstName string `json:"firstName,omitempty" db:"first_name,omitempty"`
	LastName  string `json:"lastName,omitempty" db:"last_name,omitempty"`
	Email     string `json:"email,omitempty" db:"email,omitempty"`
	Class     string `json:"class,omitempty" db:"class,omitempty"`
}
