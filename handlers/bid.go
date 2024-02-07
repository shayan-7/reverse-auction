package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Assuming you have a Bid model
type Bid struct {
	ID          uint    `json:"id" gorm:"primary_key"`
	ProductID   uint    `json:"product_id"`
	SellerID    uint    `json:"seller_id"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	IsAccepted  bool    `json:"is_accepted"`
	IsDiscarded bool    `json:"is_discarded"`
}

// @Summary Make an offer on a product
// @Description Make an offer on a product that is requested by a buyer.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param input body Bid true "Offer details"
// @Security ApiKeyAuth
// @Success 201 {object} Bid
// @Failure 400 {object} map[string]interface{}
// @Router /products/{id}/offers [post]
func makeOffer(c *gin.Context) {
	productID := c.Param("id")

	// Check if the product exists
	var product Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if the product is still open for offers
	if product.Status != Active {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not open for offers"})
		return
	}

	// Extract the seller ID from the token
	sellerID, err := extractSellerIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	var offer Bid
	if err := c.ShouldBindJSON(&offer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the product ID and seller ID for the offer
	offer.ProductID = product.ID
	offer.SellerID = sellerID

	// Create the offer
	db.Create(&offer)

	c.JSON(http.StatusCreated, offer)
}

// @Summary Get offers for a product
// @Description Get a list of offers for a specific product.
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Security ApiKeyAuth
// @Success 200 {array} Bid
// @Router /products/{id}/offers [get]
func getOffers(c *gin.Context) {
	productID := c.Param("id")

	// Check if the product exists
	var product Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	var offers []Bid
	db.Where("product_id = ?", product.ID).Find(&offers)

	c.JSON(http.StatusOK, offers)
}

// @Summary Reject an offer
// @Description Accept an offer for a product by the requester.
// @Accept json
// @Produce json
// @Param id path int true "Offer ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Router /offers/{id}/reject [put]
func rejectOffer(c *gin.Context) {
	offerID := c.Param("id")

	// Check if the user is authenticated and authorized to accept offers
	userID, err := extractSellerIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	// Fetch the offer from the database
	var offer Bid
	if err := db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
		return
	}

	// Check if the user is the requester of the product
	var product Product
	if err := db.Where("id = ?", offer.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	if product.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Mark the offer as accepted (perform your logic here)
	// For example, update the status of the offer in the database
	offer.IsAccepted = false
	db.Save(&offer)

	c.Status(http.StatusNoContent)
}

// @Summary Accept an offer
// @Description Accept an offer for a product by the requester.
// @Accept json
// @Produce json
// @Param id path int true "Offer ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} map[string]interface{}
// @Router /offers/{id}/accept [put]
func acceptOffer(c *gin.Context) {
	offerID := c.Param("id")

	// Check if the user is authenticated and authorized to accept offers
	userID, err := extractSellerIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	// Fetch the offer from the database
	var offer Bid
	if err := db.Where("id = ?", offerID).First(&offer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offer not found"})
		return
	}

	// Check if the user is the requester of the product
	var product Product
	if err := db.Where("id = ?", offer.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	log.Println(">>>>: ", product.UserID)

	if product.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		return
	}

	// Mark the offer as accepted (perform your logic here)
	// For example, update the status of the offer in the database
	offer.IsAccepted = true
	db.Save(&offer)

	c.Status(http.StatusNoContent)
}

// Helper function to extract seller ID from the token
func extractSellerIDFromToken(c *gin.Context) (uint, error) {
	claims, exists := c.Get("claims")
	if !exists {
		return 0, fmt.Errorf("token claims not found")
	}

	token, ok := claims.(*Token)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	return token.UserID, nil
}
