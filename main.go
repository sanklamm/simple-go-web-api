package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sanklamm/simple-go-web-api/database"
	"github.com/sanklamm/simple-go-web-api/models"
	"net/http"
)

func main() {
	database.ConnectDatabase()

	router := gin.Default()

	// User routes
	router.GET("/users", getUsers)
	router.POST("/users", createUser)

	// Product routes
	router.GET("/products", getProducts)
	router.POST("/products", createProduct)

	router.Run(":8080")
}

func getUsers(c *gin.Context) {
	var users []models.User
	database.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&user)
	c.JSON(http.StatusCreated, user)
}

func getProducts(c *gin.Context) {
	var products []models.Product
	database.DB.Find(&products)
	c.JSON(http.StatusOK, products)
}

func createProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&product)
	c.JSON(http.StatusCreated, product)
}
