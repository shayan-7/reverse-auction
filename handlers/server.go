package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "uniproject/docs"
)

// @title Reverse Auction API
// @version 1.0
// @description API for the reverse auction project
// @contact.name Reverse Auction Team
// @host localhost:8080
func Run() {
	// Connect to the database
	db, err = gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// AutoMigrate will attempt to automatically migrate the schema
	db.AutoMigrate(&User{}, &Product{}, &Bid{})

	// Set up the HTTP router
	router := gin.Default()

	// Define routes
	router.POST("/admin", createAdmin)
	router.POST("/signup", signUp)
	router.POST("/login", logIn)

	apiGroup := router.Group("/api")
	profileGroup := router.Group("/profile")
	profileGroup.Use(authMiddleware)
	profileGroup.GET("", userProfile)

	productAuthGroup := apiGroup.Group("")
	productAuthGroup.Use(authMiddleware)
	productAuthGroup.POST("/products", requestProduct)
	productAuthGroup.POST("/products/:id/discard", discardProductRequest)
	productAuthGroup.POST("/products/:id/approve", approveProductRequest)
	productAuthGroup.POST("/products/:id/offers", makeOffer)
	productAuthGroup.GET("/products/:id/offers", getOffers)

	productAuthGroup.POST("/offers/:id/discard", discardOffer)
	productAuthGroup.POST("/offers/:id/approve", approveOffer)

	productGroup := apiGroup.Group("")
	productGroup.GET("/products", listProducts)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
