package models

type Student struct {
	Id        int    `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	Class     string `json:"class,omitempty"`
}
