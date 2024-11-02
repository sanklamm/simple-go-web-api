package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	m.Run()
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
