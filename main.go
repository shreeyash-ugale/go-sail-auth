package main

import (
	"github.com/shreeyash-ugale/go-sail-server/controllers"
	"github.com/shreeyash-ugale/go-sail-server/database"
	"github.com/shreeyash-ugale/go-sail-server/middlewares"
	"github.com/shreeyash-ugale/go-sail-server/models"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect("host=localhost user=postgres password=9999 dbname=go-sail-auth port=5432")
	database.Migrate()
	initPlansAndActions()
	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/token", controllers.GenerateToken)
		//api.POST("/user/register", controllers.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", controllers.Ping)
		}
		api.POST("/signup", controllers.Signup)
		api.POST("/upgrade", controllers.UpgradePlan)
	}
	return router
}

func initPlansAndActions() {
	// Define actions
	templateGenerate := models.Action{Name: "Template Generate", Description: "Generate templates"}
	dockerFileGenerate := models.Action{Name: "Docker File Generate", Description: "Generate Docker files"}
	codeEvaluation := models.Action{Name: "Code Evaluation", Description: "Evaluate code"}
	securityCheck := models.Action{Name: "Security Check", Description: "Perform security checks"}

	// Define plans
	freePlan := models.Plan{Name: "Free", Description: "Free plan with basic features", Actions: []models.Action{templateGenerate, dockerFileGenerate}}
	premiumPlan := models.Plan{Name: "Premium", Description: "Premium plan with additional features", Actions: []models.Action{templateGenerate, dockerFileGenerate, codeEvaluation}}
	executivePlan := models.Plan{Name: "Executive", Description: "Executive plan with all features", Actions: []models.Action{templateGenerate, dockerFileGenerate, codeEvaluation, securityCheck}}

	// Save actions and plans to the database
	database.Instance.Create(&templateGenerate)
	database.Instance.Create(&dockerFileGenerate)
	database.Instance.Create(&codeEvaluation)
	database.Instance.Create(&securityCheck)
	database.Instance.Create(&freePlan)
	database.Instance.Create(&premiumPlan)
	database.Instance.Create(&executivePlan)
}
