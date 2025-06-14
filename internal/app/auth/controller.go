package auth

import (
	"context"
	"log"
	userpassword "microservices/authentication/internal/app/user-password"
	usertoken "microservices/authentication/internal/app/user-token"
	"microservices/authentication/internal/constants"
	"microservices/authentication/internal/database"
	"microservices/authentication/internal/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var validate = validator.New()

func HandleRegistration(c *gin.Context) {
	dbName := os.Getenv("DB_NAME")
	secretKey := os.Getenv("JWT_SECRET_KEY")

	if dbName == "" || secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate JWT tokens!"})
		return
	}

	var registeringUser models.RegisteringUserSchema
	ctx, cancel := context.WithTimeout(context.Background(), constants.ProgressTimeout)

	if err := c.ShouldBindJSON(&registeringUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		defer cancel()
		return
	}

	if err := validate.Struct(registeringUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		defer cancel()
		return
	}

	dbClient := database.GetDBInstance()
	usersCollection := database.OpenCollection(dbClient, dbName, constants.COLLECTION_NAMES["USERS"])

	checkingUsersExistingQuery := bson.M{"email": registeringUser.Email}

	defer cancel()

	if numOfUsers, err := usersCollection.CountDocuments(ctx, checkingUsersExistingQuery); err == nil && numOfUsers > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "This email has been used before!"})
		return
	}

	currentSession, err := dbClient.StartSession()
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer currentSession.EndSession(context.Background())
	userId := primitive.NewObjectID()
	userUUID := uuid.New()
	userRole := "USER"
	currentTime, _ := time.Parse(time.RFC3339, time.Now().UTC().Format(time.RFC3339))

	newUser := models.UserSchema{
		ID:   userId,
		UUID: userUUID.String(),
		Role: userRole,
		// Username:    nil,
		Email:     registeringUser.Email,
		FirstName: registeringUser.FirstName,
		LastName:  registeringUser.LastName,
		// PhoneNumber: nil,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	_, err = usersCollection.InsertOne(ctx, newUser)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = userpassword.UpdateUserPassword(ctx, newUser.UUID, registeringUser.Password)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: set user status as PENDING_VEIRIFICATION, send email to user to verify it with a token

	c.JSON(http.StatusOK, gin.H{})
}

func HandleSigningIn(c *gin.Context) {
	dbName := os.Getenv("DB_NAME")
	secretKey := os.Getenv("JWT_SECRET_KEY")

	if dbName == "" || secretKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate JWT tokens!"})
		return
	}

	var credentials models.UserSigningInSchema
	ctx, cancel := context.WithTimeout(context.Background(), constants.ProgressTimeout)

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		defer cancel()
		return
	}

	if err := validate.Struct(credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		defer cancel()
		return
	}

	if !credentials.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either email or username must be provided"})
		defer cancel()
		return
	}

	dbClient := database.GetDBInstance()
	usersCollection := database.OpenCollection(dbClient, dbName, constants.COLLECTION_NAMES["USERS"])

	defer cancel()

	var matchingUser models.UserSchema

	if err := usersCollection.FindOne(ctx, bson.M{"email": credentials.Email}).Decode(&matchingUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password is invalid"})
		return
	}

	if err := userpassword.VerifyUserPassword(matchingUser.UUID, credentials.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password is invalid"})
		return
	}

	userTokens, err := usertoken.CreateNewToken(matchingUser)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens due to some unexpected issues!"})
	}

	accessToken := userTokens.AccessToken
	refreshToken := userTokens.RefreshToken

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken, "refreshToken": refreshToken})
}
