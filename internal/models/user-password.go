package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserPasswordSchema struct {
	ID        primitive.ObjectID `bson:"_id"`
	UUID      string             `bson:"uuid"`
	UserUUID  string             `bson:"user_uuid"`
	Hash      string             `bson:"hash"`
	Salt      string             `bson:"salt"`
	Algorithm string             `bson:"algorithm"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
