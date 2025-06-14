package userpassword

import (
	"context"
	"microservices/authentication/internal/database"
	"microservices/authentication/internal/models"
	"microservices/authentication/internal/utils"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	progressTimeout             = 100 * time.Second
	userPasswordsCollectionName = "user_passwords"
)

func VerifyUserPassword(userId string, userPassword string) error {
	dbName := os.Getenv("DB_NAME")
	dbClient := database.GetDBInstance()
	userPasswordsCollection := database.OpenCollection(dbClient, dbName, userPasswordsCollectionName)

	var existingUserPasswordDetails models.UserPasswordSchema
	ctx, cancel := context.WithTimeout(context.Background(), progressTimeout)

	if err := userPasswordsCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&existingUserPasswordDetails); err != nil {
		defer cancel()
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUserPasswordDetails.Hash), []byte(userPassword)); err != nil {
		defer cancel()
		return err
	}

	defer cancel()

	return nil
}

func UpdateUserPassword(ctx context.Context, userId string, userPassword string) error {
	dbName := os.Getenv("DB_NAME")
	dbClient := database.GetDBInstance()
	userPasswordsCollection := database.OpenCollection(dbClient, dbName, userPasswordsCollectionName)

	var existingUserPasswordDetails models.UserPasswordSchema
	var newUserPasswordDetails models.UserPasswordSchema

	findingExistingPasswordQuery := bson.M{"user_id": userId}
	err := userPasswordsCollection.FindOne(ctx, findingExistingPasswordQuery).Decode(&existingUserPasswordDetails)

	currentTime, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	if err != nil {
		id := primitive.NewObjectID()
		salt, err := utils.GenerateSalt(bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		combinedPasswordAndSalt := append([]byte(userPassword), []byte(salt)...)
		hashedPassword, err := bcrypt.GenerateFromPassword(combinedPasswordAndSalt, bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		newUserPasswordDetails = models.UserPasswordSchema{
			ID:        id,
			UserID:    userId,
			Hash:      string(hashedPassword),
			Salt:      salt,
			Algorithm: "bcrypt",
			UpdatedAt: currentTime,
		}

		_, err = userPasswordsCollection.InsertOne(ctx, newUserPasswordDetails)

		return err
	}

	combinedPasswordAndSalt := append([]byte(userPassword), []byte(existingUserPasswordDetails.Salt)...)
	hashedPassword, err := bcrypt.GenerateFromPassword(combinedPasswordAndSalt, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var updatingData primitive.D
	updatingData = append(updatingData, primitive.E{Key: "hash", Value: string(hashedPassword)})
	updatingData = append(updatingData, primitive.E{Key: "updated_at", Value: currentTime})

	upsert := false
	opt := options.UpdateOne().SetUpsert(upsert)
	_, err = userPasswordsCollection.UpdateOne(
		ctx,
		findingExistingPasswordQuery,
		bson.D{{Key: "$set", Value: updatingData}},
		opt,
	)

	return err
}
