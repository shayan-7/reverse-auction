package handlers

import (
	"net/http"
	"strconv"

	_ "uniproject/docs"

	"github.com/gin-gonic/gin"
)

type Status int

const (
	Active Status = iota
	Accepted
)

type Product struct {
	ID          uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      Status `json:"status,omitempty"`
	IsDiscarded bool   `json:"is_discarded"`
	UserID      uint   `json:"user_id,omitempty"`
	User        *User  `json:"user,omitempty"`
}

// @Summary Register a new product
// @Description Register a new product by providing details such as name, description, and buyer's ID.
// @Accept json
// @Produce json
// @Param input body Product true "Product details"
// @Success 201 {object} handlers.Product
// @Router /products [post]
func requestProduct(c *gin.Context) {
	var product Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the seller ID from the token
	buyerID, err := extractSellerIDFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
		return
	}

	// Set default status to 'active'
	product.Status = Active
	product.UserID = buyerID

	// Create the product
	db.Create(&product)

	c.JSON(http.StatusCreated, product)
}

// @Summary List all products
// @Description Get a list of all products with optional sorting and filtering.
// @Accept json
// @Produce json
// @Param sort query string false "Sort field (e.g., title, price)"
// @Param filter query string false "Filter products by name"
// @Param user_id query int false "Filter products by user id"
// @Success 200 {array} Product
// @Router /products [get]
func listProducts(c *gin.Context) {
	var products []Product
	query := db

	// Sorting
	sortParam := c.Query("sort")
	if sortParam != "" {
		query = query.Order(sortParam)
	}

	// Filtering
	filterParam := c.Query("filter")
	if filterParam != "" {
		query = query.Where("title LIKE ?", "%"+filterParam+"%")
	}

	userIDParam := c.Query("user_id")
	if userIDParam != "" {
		userID, _ := strconv.Atoi(userIDParam)
		query = query.Where("user_id = ?", userID)
	}

	query.Find(&products)

	c.JSON(http.StatusOK, products)
}
