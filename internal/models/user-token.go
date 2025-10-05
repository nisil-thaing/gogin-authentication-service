package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserTokenSchema struct {
	ID           primitive.ObjectID `bson:"_id"`
	UUID         string             `bson:"uuid"`
	UserID       string             `bson:"user_id"`
	RefreshToken string             `bson:"refresh_token"`
	ExpiresAt    time.Time          `bson:"expires_at"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
	RevokedAt    time.Time          `bson:"revoked_at,omitempty"`
}

type UserTokensPublicInfo struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type JWTSigningClaims struct {
	UserID    string // replace it by something like UserUUID once it available
	TokenUUID string `json:"jti"`
	Username  string
	Email     string
	FirstName string
	LastName  string
	Role      string
	jwt.RegisteredClaims
}
