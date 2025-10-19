package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	TokenVerificationType string
)

var TokenVerificationTypes = map[string]string{
	"EMAIL_VERIFICATION":        "Verify user's email during registration or email change",
	"PASSWORD_RESET":            "Reset user password after a forgot-password request",
	"TWO_FACTOR":                "Verify user identity in two-factor authentication (2FA/MFA)",
	"EMAIL_CHANGE_VERIFICATION": "Confirm user's new email address before applying the change",
	"ACCOUNT_REACTIVATION":      "Reactivate a suspended or deactivated user account",
	"INVITE":                    "Invite a user to join a team, workspace, or project",
	"MAGIC_LOGIN":               "Passwordless login via magic link sent to user's email",
	"DEVICE_VERIFICATION":       "Verify a new device or login attempt from an unknown location",
	"BILLING_APPROVAL":          "Confirm a billing or subscription-related action",
}

type UserVerificationSchema struct {
	ID        primitive.ObjectID    `bson:"_id"`
	UUID      string                `bson:"uuid"`
	UserUUID  string                `bson:"user_uuid" validate:"required"`
	Token     string                `bson:"token" validate:"required"`
	Type      TokenVerificationType `bson:"type" validate:"required, eq=EMAIL_VERIFICATION"` // TODO: update it if necessary
	CreatedAt time.Time             `bson:"created_at"`
	ExpiresAt time.Time             `bson:"expires_at" validate:"required"`
}

type UserVerificationPublicInfo struct {
	Token     string
	ExpiresAt time.Time
}

const (
	TokenVerificationEmail               TokenVerificationType = "EMAIL_VERIFICATION"
	TokenVerificationPasswordReset       TokenVerificationType = "PASSWORD_RESET"
	TokenVerificationTwoFactor           TokenVerificationType = "TWO_FACTOR"
	TokenVerificationEmailChange         TokenVerificationType = "EMAIL_CHANGE_VERIFICATION"
	TokenVerificationAccountReactivation TokenVerificationType = "ACCOUNT_REACTIVATION"
	TokenVerificationInvite              TokenVerificationType = "INVITE"
	TokenVerificationMagicLogin          TokenVerificationType = "MAGIC_LOGIN"
	TokenVerificationDevice              TokenVerificationType = "DEVICE_VERIFICATION"
	TokenVerificationBillingApproval     TokenVerificationType = "BILLING_APPROVAL"
)
