package main

import (
	"seamlessop/controllers"
	"seamlessop/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// r.GET("/books/:id", controllers.FindBook)
	// r.POST("/books", controllers.CreateBook)

	api := r.Group("/seamless/api")
	api.Use(models.ConnectToDB())
	{
		api.GET("/getToken", controllers.UserCreateToken)
		api.GET("/validateToken", controllers.UserValidToken)
		api.POST("/reserve", controllers.BetReserve)
		api.GET("/cancelReserve", controllers.BetCancelReserve)
		api.POST("/debitReserve", controllers.BetDebitReserve)
		api.GET("/commitReserve", controllers.BetCommitReserve)
		api.POST("/debitCustomer", controllers.BetDebitCustomer)
		api.POST("/creditCustomer", controllers.BetCreditCustomer)
	}

	r.Run()
}
