package usertoken

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

func CreateNewToken(user models.UserSchema) (*models.UserTokensPublicInfo, error) {
	dbName := os.Getenv("DB_NAME")
	secretKey := os.Getenv("JWT_SECRET_KEY")

	dbClient := database.GetDBInstance()
	userTokensCollection := database.OpenCollection(dbClient, dbName, constants.COLLECTION_NAMES["USER_TOKENS"])

	ctx, cancel := context.WithTimeout(context.Background(), constants.ProgressTimeout)

	tokenUUID := uuid.New()
	userTokens, err := utils.GenerateTokens(user, tokenUUID.String(), secretKey)
	defer cancel()

	if err != nil {
		return nil, err
	}

	if userTokens == nil {
		return nil, errors.New("❌ Could not generate the tokens")
	}

	tokenId := primitive.NewObjectID()
	currentTime, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))

	newTokenData := &models.UserTokenSchema{
		ID:           tokenId,
		UUID:         tokenUUID.String(),
		UserUUID:     user.UUID,
		RefreshToken: userTokens.RefreshToken,
		ExpiresAt:    userTokens.ExpiresAt,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}

	_, err = userTokensCollection.InsertOne(ctx, newTokenData)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return userTokens, nil
}

func UpdateAccessToken() {}

func RevokeToken() {}
