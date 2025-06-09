package response

type User struct {
	Model
	Name    string `json:"name" example:"John"`
	Surname string `json:"surname" example:"Doe"`
}

type UserResponse struct {
	OK   bool `json:"ok" example:"true"`
	Data User `json:"data,omitempty"`
}
