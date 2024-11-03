package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sanklamm/simple-go-web-api/database"
	"github.com/sanklamm/simple-go-web-api/models"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	database.TestDB = true
	database.ConnectDatabase()
	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	user := models.User{Name: "New User", Email: "newuser@example.com"}
	jsonValue, _ := json.Marshal(user)

	// create a request
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// perform the request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, resp.Code)

	var createdUser models.User
	err := json.Unmarshal(resp.Body.Bytes(), &createdUser)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)
	assert.NotEqual(t, uuid.Nil, createdUser.ID)
}

func TestGetUsers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	// create a test user
	user := models.User{Name: "Test User", Email: "test@example.com"}
	database.DB.Create(&user)

	// create a request to send to the endpoint
	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	resp := httptest.NewRecorder()

	// perform the request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)

	var users []models.User
	err := json.Unmarshal(resp.Body.Bytes(), &users)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(users), 1)
}

func TestGetUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	// create a test user
	user := models.User{Name: "Test User", Email: "testuser@example.com"}
	database.DB.Create(&user)

	// create a request to get the user by ID
	req, _ := http.NewRequest(http.MethodGet, "/users/"+user.ID.String(), nil)
	resp := httptest.NewRecorder()

	// perform the request
	router.ServeHTTP(resp, req)

	// Assertions
	var returnedUser models.User
	err := json.Unmarshal(resp.Body.Bytes(), &returnedUser)
	assert.Nil(t, err)
	assert.Equal(t, user.Name, returnedUser.Name)
}

func TestGetUserInvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := setupRouter()

	// create a request to get the user by ID
	req, _ := http.NewRequest(http.MethodGet, "/users/invalid-uuid", nil)
	resp := httptest.NewRecorder()

	// perform the request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}
