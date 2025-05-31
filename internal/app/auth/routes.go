package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.RouterGroup) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is Auth Routes"})
	})

	router.POST("/register", HandleRegistration)
	// router.POST("/signin", HandleSigningIn)
}
