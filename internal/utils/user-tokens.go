package utils

import (
	"errors"
	"log"
	"microservices/authentication/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateUserTokens(user models.UserSchema, tokenUUID string, secretKey string) (*models.UserTokensPublicInfo, error) {
	timeNow := time.Now().UTC()
	accessTokenClaims := models.JWTSigningClaims{
		UserUUID:  user.UUID,
		TokenUUID: tokenUUID,
		Email:     user.Email,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(time.Duration(24) * time.Hour)),
		},
	}

	tokensExpireAt := timeNow.Add(time.Duration(168) * time.Hour)
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

	return &models.UserTokensPublicInfo{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
		ExpiresAt:    tokensExpireAt,
	}, nil
}

func ValidateToken(token string, secretKey string) (*models.JWTSigningClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &models.JWTSigningClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	return parsedToken.Claims.(*models.JWTSigningClaims), nil
}
