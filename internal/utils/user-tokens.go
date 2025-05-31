package utils

import (
	"log"
	"microservices/authentication/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(user models.UserSchema, secretKey string) (*models.UserTokensPublicInfo, error) {
	timeNowInLocal := time.Now().Local()
	accessTokenClaims := models.JWTSigningClaims{
		UserID:    user.UserID,
		Email:     user.Email,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeNowInLocal.Add(time.Duration(24) * time.Hour)),
		},
	}

	tokensExpireAt := timeNowInLocal.Add(time.Duration(168) * time.Hour)
	refreshTokenClaims := models.JWTSigningClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokensExpireAt),
		},
	}

	signedAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(secretKey))
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	signedRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(secretKey))
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	expiresAt, _ := time.Parse(time.RFC3339, tokensExpireAt.UTC().Format(time.RFC3339))
	return &models.UserTokensPublicInfo{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}
