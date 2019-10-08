package main

import (
	"github.com/gin-gonic/gin"
	"github.com/theMomax/GolangPizzaServiceManager/controllers"
)

func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1/")
	{
		v1.POST("/order/:id", controllers.Order)
		v1.GET("/store", controllers.FetchAvailable)
		v1.POST("/refill", controllers.AddResource)
		v1.PUT("/recipe", controllers.CreateRecipe)
		v1.GET("/recipe/:id", controllers.Fetch)
		v1.GET("/recipe", controllers.FetchAll)
		v1.DELETE("/recipe/:id", controllers.DeleteRecipe)
		v1.PUT("/recipe/:id", controllers.UpdateRecipe)

		v1.GET("/price/:id", controllers.Price)
	}
	router.Run()
}
