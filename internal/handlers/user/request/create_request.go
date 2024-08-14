package request

type Request struct {
	FirstName string `json:"firstName" validate:"required,min=2,max=50,alphanum"`
	LastName  string `json:"lastName" validate:"required,min=2,max=50,alphanum"`
}
