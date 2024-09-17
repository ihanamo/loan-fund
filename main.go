package main

import (
	"loan-fund/database"
	"loan-fund/handlers"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	database.InitDB()

	e.POST("/User/RegisterUser", handlers.CreateUser)
	e.POST("/User/LoginUser", handlers.LoginUser)

}