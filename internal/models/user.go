package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UserRole   string
	UserStatus string
)

var UserRoles = map[string]string{
	"ADMIN": "User with administrative privileges and full access",
	"USER":  "Regular user with standard access rights",
}

var UserStatuses = map[string]string{
	"PENDING_VERIFICATION": "User has registered but not yet verified their email",
	"INACTIVE":             "User account is temporarily disabled or suspended",
	"ACTIVE":               "User account is active and fully verified",
}

type UserSchema struct {
	ID          primitive.ObjectID `bson:"_id"`
	UUID        string             `bson:"uuid"`
	Role        UserRole           `bson:"role" validate:"required,eq=ADMIN|eq=USER"`
	Username    *string            `bson:"username,omitempty"`
	Email       string             `bson:"email" validate:"email,required"`
	FirstName   *string            `bson:"first_name" validate:"min=2,max=100"`
	LastName    *string            `bson:"last_name" validate:"min=2,max=100"`
	PhoneNumber *string            `bson:"phone_number,omitempty" validate:"min=10"`
	AvatarUrl   *string            `bson:"avatar_url,omitempty"`
	Status      UserStatus         `bson:"status" validate:"required,eq=PENDING_VERIFICATION|eq=INACTIVE|ACTIVE"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at,omitempty"`
}

const (
	UserAdminRole  UserRole = "ADMIN"
	UserNormalRole UserRole = "USER"
)

const (
	UserPendingVerificationStatus UserStatus = "PENDING_VERIFICATION"
	UserInactiveStatus            UserStatus = "INACTIVE"
	UserActiveStatus              UserStatus = "ACTIVE"
)
