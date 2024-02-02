package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Assuming you have a Bid model
type Bid struct {
	ID          uint    `gorm:"primary_key"`
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
