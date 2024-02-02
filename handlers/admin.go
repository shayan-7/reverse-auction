package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Register a new user
// @Description Register a new user by providing a unique username and password.
// @Accept json
// @Produce json
// @Param input body User true "User registration details"
// @Success 201 {object} map[string]interface{}
// @Router /admin [post]
func createAdmin(c *gin.Context) {
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

	user.IsAdmin = true
	// Hash the password before saving it to the database (you should use a secure hashing library)
	// For simplicity, we are not doing password hashing in this example
	db.Create(&user)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// @Summary Discard a product request
// @Description Discard a product request by the admin user.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Router /products/{id}/discard [delete]
func discardProductRequest(c *gin.Context) {
	productID := c.Param("id")

	// Check if the user is an admin
	isAdmin, err := isAdminUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Discard the product request (perform your logic here)
	var product Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Set is_discarded to true
	product.IsDiscarded = true
	db.Save(&product)

	c.Status(http.StatusNoContent)
}

// @Summary Discard a product request
// @Description Discard a product request by the admin user.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Router /products/{id}/approve [delete]
func approveProductRequest(c *gin.Context) {
	productID := c.Param("id")

	// Check if the user is an admin
	isAdmin, err := isAdminUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Discard the product request (perform your logic here)
	var product Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Set is_discarded to true
	product.IsDiscarded = false
	db.Save(&product)

	c.Status(http.StatusNoContent)
}

// @Summary Discard an offer
// @Description Discard an offer by the admin user.
// @Accept json
// @Produce json
// @Param id path int true "Offer ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Router /offers/{id}/approve [delete]
func approveOffer(c *gin.Context) {
	offerID := c.Param("id")

	// Check if the user is an admin
	isAdmin, err := isAdminUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Discard the offer (perform your logic here)
	var bid Bid
	if err := db.Where("id = ?", offerID).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bid not found"})
		return
	}

	// Set is_discarded to true
	bid.IsDiscarded = false
	db.Save(&bid)

	c.Status(http.StatusNoContent)
}

// @Summary Discard an offer
// @Description Discard an offer by the admin user.
// @Accept json
// @Produce json
// @Param id path int true "Offer ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Router /offers/{id}/discard [delete]
func discardOffer(c *gin.Context) {
	offerID := c.Param("id")

	// Check if the user is an admin
	isAdmin, err := isAdminUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Discard the offer (perform your logic here)
	var bid Bid
	if err := db.Where("id = ?", offerID).First(&bid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bid not found"})
		return
	}

	// Set is_discarded to true
	bid.IsDiscarded = true
	db.Save(&bid)

	c.Status(http.StatusNoContent)
}

// Helper function to check if the user is an admin
func isAdminUser(c *gin.Context) (bool, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return false, fmt.Errorf("token claims not found")
	}

	token, ok := claims.(*Token)
	if !ok {
		return false, fmt.Errorf("invalid token claims")
	}

	return token.IsAdmin, nil
}
