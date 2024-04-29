package request

type SignUpRequest struct {
	Firstname string `json:"firstname" validate:"min=2"`
	Lastname  string `json:"lastname" validate:"min=2"`
	Email     string `json:"email" validate:"email"`
	Password  string `json:"password" validate:"-"`
}
