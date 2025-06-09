package models

type RegisteringUserSchema struct {
	Email           string  `bson:"email" validate:"email,required"`
	FirstName       *string `bson:"first_name" validate:"min=2,max=100"`
	LastName        *string `bson:"last_name" validate:"min=2,max=100"`
	Password        string  `bson:"password" validate:"required,min=8,eqfield=PasswordConfirm"`
	PasswordConfirm string  `bson:"password_confirm" validate:"required"`
}

// type UserSigningInSchema struct {
// 	Email    string `bson:"email" validate:"email,required"`
// 	Password string `bson:"password" validate:"required"`
// }
