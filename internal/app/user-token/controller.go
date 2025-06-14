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
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	progressTimeout          = 100 * time.Second
	userTokensCollectionName = "user_tokens"
)

func UpdateUserToken(user models.UserSchema, secretKey string) (*models.UserTokensPublicInfo, error) {
	dbName := os.Getenv("DB_NAME")
	dbClient := database.GetDBInstance()
	userTokensCollection := database.OpenCollection(dbClient, dbName, userTokensCollectionName)

	var userTokensDetails models.UserTokenSchema
	ctx, cancel := context.WithTimeout(context.Background(), progressTimeout)

	findingExistingTokenQuery := bson.M{"user_id": user.UserID}

	err := userTokensCollection.FindOne(ctx, findingExistingTokenQuery).Decode(&userTokensDetails)

	var updatingData bson.D
	currentTime, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))

	isNoDataExisting := err == mongo.ErrNoDocuments

	if err != nil && !isNoDataExisting {
		defer cancel()
		return nil, err
	}

	isTokensExpired := userTokensDetails.ExpiresAt.Before(currentTime)

	if !isTokensExpired {
		_, err = utils.ValidateToken(userTokensDetails.AccessToken, secretKey)
		isTokensExpired = err != nil
	}

	if isNoDataExisting || isTokensExpired {
		// TODO: No existing user's tokens found, or it's found, but the tokens are expired:
		// generate new tokens, add to DB"
		userTokens, err := utils.GenerateTokens(user, secretKey)
		defer cancel()

		if err != nil {
			return nil, err
		}

		if userTokens == nil {
			return nil, errors.New("❌ Could not generate the tokens")
		}

		if isNoDataExisting {
			updatingData = append(updatingData, bson.E{Key: "_id", Value: primitive.NewObjectID()})
			updatingData = append(updatingData, bson.E{Key: "user_id", Value: user.UserID})
			updatingData = append(updatingData, bson.E{Key: "created_at", Value: currentTime})
		}

		updatingData = append(updatingData, bson.E{Key: "access_token", Value: userTokens.AccessToken})
		updatingData = append(updatingData, bson.E{Key: "refresh_token", Value: userTokens.RefreshToken})
		updatingData = append(updatingData, bson.E{Key: "expires_at", Value: userTokens.ExpiresAt})
		updatingData = append(updatingData, bson.E{Key: "updated_at", Value: currentTime})

		// Store user token to database
		upsert := true
		opts := options.UpdateOne().SetUpsert(upsert)

		_, err = userTokensCollection.UpdateOne(ctx, findingExistingTokenQuery, bson.M{
			"$set": updatingData,
		}, opts)
		if err != nil {
			return nil, err
		}

		return userTokens, nil
	}

	defer cancel()

	return &models.UserTokensPublicInfo{
		AccessToken:  userTokensDetails.AccessToken,
		RefreshToken: userTokensDetails.RefreshToken,
		ExpiresAt:    userTokensDetails.ExpiresAt,
	}, nil
}
