package entity

type User struct {
	Model   Model  `json:"model"`
	Name    string `json:"name,omitempty"`
	Surname string `json:"surname,omitempty"`
}
