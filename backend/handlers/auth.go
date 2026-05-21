package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Login(c *gin.Context) {

	var user struct {
		UserEmail string `json:"useremail"`
		Password  string `json:"password"`
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.UserEmail = strings.TrimSpace(user.UserEmail)
	if user.UserEmail == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "useremail and password are required"})
		return
	}

	var loginUser models.User
	if err := db.DB.Where("email = ?", user.UserEmail).First(&loginUser).Error; err != nil {
		log.Printf("User not found %s. Error from loginHandler...", user.UserEmail)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return

	}

	if loginUser.PasswordHash != user.Password {
		log.Printf("Invalid password for user %s", user.UserEmail)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Printf("jwt secret not found in environment , check for secret config ")
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "ACCESS TOKEN NOT CONFIGURED , INTERNAL SERVER SIDE ERROR ."})
		return
	}

	// Generate token. Store `user_id` as a string UUID to keep token claim types stable.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_email": loginUser.Email,
		"user_id":    loginUser.ID.String(),
		//setting expiry to 10000 hours for testing , to be changed in prod.
		"exp": time.Now().Add(time.Hour * 10000).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Failed to sign token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log.Printf("Token generated successfully for user: %s (ID: %s)", loginUser.Email, loginUser.ID.String())

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"userEmail": loginUser.Email,
			//	"userId":    loginUser.ID,
			//"phone":      loginUser.Phoneno,
			//"userId": loginUser.ID.String(),
		},
	})
}
