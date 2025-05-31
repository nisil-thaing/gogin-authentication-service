package usertoken

import (
	"context"
	"errors"
	"microservices/authentication/internal/database"
	"microservices/authentication/internal/models"
	"microservices/authentication/internal/utils"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const userTokensCollectionName = "user_tokens"

func UpdateUserToken(ctx context.Context, user models.UserSchema, secretKey string) (*models.UserTokensPublicInfo, error) {
	dbName := os.Getenv("DB_NAME")
	dbClient := database.GetDBInstance()
	userTokensCollection := database.OpenCollection(dbClient, dbName, userTokensCollectionName)

	var userTokensDetails models.UserTokenSchema

	findingExistingTokenQuery := bson.M{"user_id": user.UserID}

	err := userTokensCollection.FindOne(ctx, findingExistingTokenQuery).Decode(&userTokensDetails)

	currentTime, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))

	// Handle the case of err == nil, userTokensDetails != nil, but the RefreshToken is expired
	if err != nil {
		// TODO: No existing user's tokens found, generate new tokens, add to DB"
		userTokens, err := utils.GenerateTokens(user, secretKey)
		if err != nil {
			return nil, err
		}

		if userTokens == nil {
			return nil, errors.New("❌ Could not generate the tokens")
		}

		userTokensDetails.ID = primitive.NewObjectID()
		userTokensDetails.UserID = user.UserID
		userTokensDetails.AccessToken = userTokens.AccessToken
		userTokensDetails.RefreshToken = userTokens.RefreshToken
		userTokensDetails.ExpiresAt = userTokens.ExpiresAt
		userTokensDetails.CreatedAt = currentTime
		userTokensDetails.UpdatedAt = currentTime

		// TODO: store it to the database
		_, err = userTokensCollection.InsertOne(ctx, userTokensDetails)
		if err != nil {
			return nil, err
		}

		return userTokens, nil
	}

	// TODO: check if the tokens has expired?
	// YES? Generate new tokens, update the existing
	// NO? renew the AccessToken by the RefreshToken, store it again

	// TODO
	return nil, nil
}
