package models

type User struct {
	ID       int    `json:"id_user"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}