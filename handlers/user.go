package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "uniproject/docs"
)

var (
	db  *gorm.DB
	err error
)

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

type Token struct {
	UserID  uint
	IsAdmin bool
	jwt.StandardClaims
}

// @Summary Register a new user
// @Description Register a new user by providing a unique username and password.
// @Accept json
// @Produce json
// @Param input body User true "User registration details"
// @Success 201 {object} map[string]interface{}
// @Router /signup [post]
func signUp(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the username already exists
	if isUsernameTaken(user.Username) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	// Hash the password before saving it to the database (you should use a secure hashing library)
	// For simplicity, we are not doing password hashing in this example
	db.Create(&user)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// @Summary Log in as a user
// @Description Log in by providing a username and password.
// @Accept json
// @Produce json
// @Param input body User true "User login details"
// @Success 200 {object} map[string]interface{}
// @Router /login [post]
func logIn(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the username and password match
	user, err := getUserByCredentials(user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	token, err := generateToken(user.ID, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Helper functions
func isUsernameTaken(username string) bool {
	var user User
	db.Where("username = ?", username).First(&user)
	return user.ID != 0
}

func getUserByCredentials(username, password string) (User, error) {
	var user User
	err := db.Where("username = ? AND password = ?", username, password).First(&user).Error
	return user, err
}

func generateToken(userID uint, isAdmin bool) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours
	claims := &Token{
		UserID:  userID,
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key")) // Replace with your secret key
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func userProfile(c *gin.Context) {
	// Access user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}

	// Retrieve user from the database using the user ID
	var user User
	db.First(&user, userID.(uint))

	c.JSON(http.StatusOK, gin.H{"username": user.Username, "userID": user.ID})
}

func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil // Replace with your secret key
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	claims, ok := token.Claims.(*Token)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Set user ID in context for handlers to access
	c.Set("claims", claims)
}
