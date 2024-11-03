package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sanklamm/simple-go-web-api/database"
	"github.com/sanklamm/simple-go-web-api/middleware"
	"github.com/sanklamm/simple-go-web-api/models"
)

func main() {
	database.ConnectDatabase()

	router := setupRouter()

	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.AuthMiddleware())

	router.POST("/login", login)

	router.GET("/", home)

	router.GET("/users", getUsers)
	router.GET("/users/:id", getUserById)
	router.POST("/users", createUser)

	router.GET("/products", getProducts)
	router.GET("/products/:id", getProductById)
	router.POST("/products", createProduct)

	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")

	return router
}

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Simplified user authentication
	var dbUser models.User
	if err := database.DB.Where("email = ?", user.Email).First(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Normally, you should check the password here

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
	})

	tokenString, err := token.SignedString(middleware.SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// GET /users
// GET /users?name=John
func getUsers(c *gin.Context) {
	var users []models.User
	name := c.Query("name")

	query := database.DB
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	query.Find(&users)
	if c.GetHeader("Hx-Request") != "" {
		c.HTML(http.StatusOK, "users.html", gin.H{"users": users})
	} else {
		c.JSON(http.StatusOK, users)
	}
}

// GET /users/:id
func getUserById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// POST /users
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
	name := c.Query("name")

	query := database.DB
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	query.Find(&products)
	c.JSON(http.StatusOK, products)
}

func getProductById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
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
