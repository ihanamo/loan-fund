package handlers

import (
	"loan-fund/database"
	"loan-fund/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

var JWTKey = []byte("secret_key")

func GenerateJWT(customer models.User) (string, error) {
	claims := &jwt.MapClaims{
		"customer-id": customer.ID,
		"username":    customer.Username,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return err
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"messgae": "Failed to hash password"})
	}
	user.Password = string(hashPass)
	log.Println("the hash password is:", user.Password)

	result := database.DB.Create(user)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, result.Error)
	}
	log.Println("user created")

	token, err := GenerateJWT(*user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to generate token"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "User created successfuly",
		"user":    user,
		"token":   token,
	})
}

func ReadCustomer(c echo.Context) error {
	usererID := c.Param("id")
	var user models.User
	if result := database.DB.First(&user, usererID); result.Error != nil {
		return c.JSON(http.StatusNotFound, result.Error)
	}

	return c.JSON(http.StatusOK, user)
}
