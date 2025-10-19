package userverification

import (
	"context"
	"errors"
	"log"
	"microservices/authentication/internal/constants"
	"microservices/authentication/internal/database"
	"microservices/authentication/internal/models"
	"microservices/authentication/internal/utils"
	"os"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateNewToken(user models.UserSchema) (*models.UserVerificationPublicInfo, error) {
	dbName := os.Getenv("DB_NAME")
	dbClient := database.GetDBInstance()
	userVerificationCollection := database.OpenCollection(dbClient, dbName, constants.COLLECTION_NAMES["USER_VERIFICATION"])

	ctx, cancel := context.WithTimeout(context.Background(), constants.ProgressTimeout)

	tokenInfo, err := utils.GenerateVerificationToken()
	defer cancel()

	if err != nil {
		return nil, err
	}

	if tokenInfo == nil || tokenInfo.Token == "" {
		return nil, errors.New("❌ Could not generate the token")
	}

	tokenId := primitive.NewObjectID()
	tokenUUID := uuid.New()
	tokenType := models.TokenVerificationEmail
	currentTime, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))

	newTokenData := &models.UserVerificationSchema{
		ID:        tokenId,
		UUID:      tokenUUID.String(),
		UserUUID:  user.UUID,
		Token:     tokenInfo.Token,
		Type:      tokenType,
		CreatedAt: currentTime,
		ExpiresAt: tokenInfo.ExpiresAt,
	}

	_, err = userVerificationCollection.InsertOne(ctx, newTokenData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &models.UserVerificationPublicInfo{
		Token:     newTokenData.Token,
		ExpiresAt: newTokenData.ExpiresAt,
	}, nil
}
