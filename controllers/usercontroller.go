package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/shreeyash-ugale/go-sail-server/database"
	"github.com/shreeyash-ugale/go-sail-server/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	var user models.User
	if c.ContentType() != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request content type must be application/json"})
		return
	}

	var reqBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Plan     string `json:"plan"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var plan models.Plan
	if err := database.Instance.Where("name = ?", reqBody.Plan).First(&plan).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Plan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = reqBody.Password
	user.Email = reqBody.Email
	user.Username = reqBody.Username
	user.PlanID = plan.ID

	// Save the user to the database
	database.Instance.Create(&user)

	// Generate an API key for the user
	apiKey := models.APIKey{
		Key:    uuid.New().String(),
		UserID: user.ID,
	}
	database.Instance.Create(&apiKey)

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "api_key": apiKey.Key})
}

func UpgradePlan(c *gin.Context) {
	if c.ContentType() != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request content type must be application/json"})
		return
	}

	var request struct {
		Email    string `json:"email"`
		PlanName string `json:"plan_name"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user by email
	var user models.User
	database.Instance.Where("email = ?", request.Email).First(&user)

	// Find the new plan by name
	var newPlan models.Plan
	database.Instance.Where("name = ?", request.PlanName).First(&newPlan)

	// Update the user's plan
	user.PlanID = newPlan.ID
	database.Instance.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Plan upgraded successfully"})
}
