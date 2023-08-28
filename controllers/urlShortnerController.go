package controllers

import (
	"errors"
	"example.com/m/initializers"
	"example.com/m/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	initializers.LoadEnviromentalVariables()
	initializers.ConnectToDatabse()
}

func GetCurrentUser(c *gin.Context) (models.User, error) {
	var user models.User
	userID, exists := c.Get("user_id")
	if !exists {
		return user, fmt.Errorf("No user ID found in the context")
	}

	if err := initializers.DB.First(&user, userID).Error; err != nil {
		return user, fmt.Errorf("Failed to get user details")
	}

	return user, nil
}

func generateShortcode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())
	var shortcode strings.Builder
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		shortcode.WriteByte(charset[randomIndex])
	}
	return shortcode.String()
}

func GetLongUrl(c *gin.Context) {
	fmt.Println("53454444", c)
	var input struct {
		LongUrl string `form:"long_url"`
	}
	var existingUrlMapping models.URlMapping
	var existingUser models.User
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	currentUser, err := GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err" +
			"or": err.Error()})
		return
	}
	if input.LongUrl == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "In valid Input",
		})
		return
	}

	result := initializers.DB.Where("long_url = ? AND user_id = ?", input.LongUrl, currentUser.ID).First(&existingUrlMapping)
	result_user := initializers.DB.Where("id = ?", currentUser.ID).First(&existingUser)
	if result_user.Error != nil {
		c.JSON(400, gin.H{
			"error": "User not Found",
		})
		return
	}
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) || result_user.Error == nil {
			shortCode := generateShortcode(8)
			newUrlMapping := models.URlMapping{LongUrl: input.LongUrl, ShortCode: shortCode, ShortUrl: "http://localhost:3000/" + shortCode, UserID: currentUser.ID}
			createResult := initializers.DB.Create(&newUrlMapping)
			if createResult.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Error in creating URL mapping",
					"details": createResult.Error.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"shortUrl": "http://localhost:3000/" + shortCode,
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Database error",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"shortUrl": "http://localhost:3000/" + existingUrlMapping.ShortCode,
		})
		return
	}

}

func RedirectingToOrignalUrl(c *gin.Context) {
	shortCode := c.Param("shortcode")
	var urlMaping models.URlMapping

	initializers.DB.Where("short_code = ? ", shortCode).First(&urlMaping)
	if urlMaping.ID == 0 {
		c.JSON(http.StatusNotFound,
			gin.H{
				"error": "record not found",
			})
	}
	fmt.Println(urlMaping.LongUrl)
	c.Redirect(http.StatusMovedPermanently, urlMaping.LongUrl)
}

func GetAllUrls(c *gin.Context) {
	var urls []models.URlMapping
	currentUser, err := GetCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err" +
			"or": err.Error()})
		return
	}
	all_urls := initializers.DB.Preload("User").Where("user_id = ? ", currentUser.ID).Find(&urls)
	if all_urls.Error != nil {
		c.JSON(200, gin.H{
			"error": "Error in fetching Data",
		})
	}
	c.JSON(200, gin.H{
		"urls": urls,
	})

}

func Signup(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "In valid Username",
		})
		return
	}
	if len(user.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Length of password must be greater than 8",
		})
		return
	}

	var existingUser models.User
	result := initializers.DB.Where("username = ?", user.Username).First(&existingUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
				return
			}

			user.Password = string(hashedPassword)
			initializers.DB.Create(&user)

			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

}

func Login(c *gin.Context) {
	var user struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var foundUser models.User
	initializers.DB.Where("username = ?", user.Username).First(&foundUser)

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	claims := models.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
		UserID: foundUser.ID,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
