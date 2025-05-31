package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPasswordSchema struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"user_id"`
	Hash      string             `bson:"hash"`
	Salt      string             `bson:"salt"`
	Algorithm string             `bson:"algorithm"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
