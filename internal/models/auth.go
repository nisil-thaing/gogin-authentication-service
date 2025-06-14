package models

type RegisteringUserSchema struct {
	Email           string  `bson:"email" validate:"email,required"`
	FirstName       *string `bson:"first_name" validate:"min=2,max=100"`
	LastName        *string `bson:"last_name" validate:"min=2,max=100"`
	Password        string  `bson:"password" validate:"required,min=8,eqfield=PasswordConfirm"`
	PasswordConfirm string  `bson:"password_confirm" validate:"required"`
}

type UserSigningInSchema struct {
	Email    *string `json:"email,omitempty" validate:"omitempty,email"`
	Username *string `json:"username,omitempty" validate:"omitempty"`
	Password string  `json:"password" validate:"required,min=8"`
}

// Validate ensures either Email or Username is provided
func (s *UserSigningInSchema) Validate() bool {
	return (s.Email != nil && *s.Email != "") || (s.Username != nil && *s.Username != "")
}
