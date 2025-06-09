package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSchema struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      string             `bson:"user_id"`
	Role        string             `bson:"role" validate:"required,eq=ADMIN|eq=USER"`
	Username    *string            `bson:"username"`
	Email       string             `bson:"email" validate:"email,required"`
	FirstName   *string            `bson:"first_name" validate:"min=2,max=100"`
	LastName    *string            `bson:"last_name" validate:"min=2,max=100"`
	PhoneNumber *string            `bson:"phone_number" validate:"min=10"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
