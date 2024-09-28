package main

import (
	"loan-fund/MiddleWare"
	"loan-fund/database"
	"loan-fund/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	database.InitDB()

	e.POST("/User/RegisterUser", handlers.CreateUser)
	e.POST("/User/LoginUser", handlers.LoginUser)

	r := e.Group("")
	r.Use(MiddleWare.JWTMiddleware())
	r.Use(MiddleWare.ExtractClaims)
	
	r.GET("/User/UserInfo", handlers.ReadUser)
	r.PUT("/User/UpdateInfo", handlers.UpdateUser)
	r.DELETE("/User/DeleteUser", handlers.DeleteUser)

	r.POST("/Transaction/Deposit", handlers.Deposit)
	r.POST("/Transaction/Loan", handlers.IssueLaon)
	r.POST("/Transaction/Repay", handlers.MakeRepayment)
	r.GET("/Transaction/History", handlers.GetTransactionHistory)


	e.Logger.Fatal(e.Start(":8080"))

}