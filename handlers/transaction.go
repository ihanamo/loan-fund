package handlers

import (
	"loan-fund/database"
	"loan-fund/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

var LoanOptions = map[string]time.Duration{
	"10month": time.Hour * 24 * 300,
	"15month": time.Hour * 24 * 450,
	"20month": time.Hour * 24 * 600,
}

func IssueLaon(c echo.Context) error {
	userIDstr := c.Param("id")

	userID, err := strconv.ParseUint(userIDstr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid user ID"})
	}

	loan := new(models.Loan)
	if err := c.Bind(loan); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request data"})
	}

	dueDateOption := c.QueryParam("due_date_option")
	if dueDateOption == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Due date option is required"})
	}

	duration, ok := LoanOptions[dueDateOption]
	if !ok {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid due date option"})
	}

	loan.UserID = uint(userID)
	loan.Status = "acive"
	loan.CreatedAt = time.Now()
	loan.DueDate = time.Now().Add(duration)

	if result := database.DB.Create(loan); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to issue loan"})
	}

	transaction := models.Transaction{
		UserID:    loan.UserID,
		LoanID:    loan.ID,
		Type:      "loan",
		Amount:    loan.Amount,
		Balance:   loan.Amount,
		CreatedAt: time.Now(),
	}

	if result := database.DB.Create(&transaction); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to log transaction"})
	}

	logEntry := models.Log{
		UserID:        loan.UserID,
		TransactionID: transaction.ID,
		Type:          "loan",
		CreatedAt:     time.Now(),
	}

	if result := database.DB.Create(&logEntry); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to log action"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message":  "Loan issued successfuly",
		"loan":     loan,
		"due_date": loan.DueDate,
		"balance":  transaction.Balance,
		"log":      logEntry,
	})
}

func MakeRepayment(c echo.Context) error {
	userID := c.Param("id")

	repayment := new(models.Transaction)
	if err := c.Bind(repayment); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid request data"})
	}

	var loan models.Loan
	if result := database.DB.Where("user_id = ? AND status = ?", userID, "active").First(&loan); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "No active loan found"})
	}

	loan.Amount -= repayment.Amount

	if loan.Amount <= 0 {
		loan.Status = "paid"
	}

	if result := database.DB.Save(&loan); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update loan"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Repayment successful",
		"balance": loan.Amount,
	})
}

func Deposit(c echo.Context) error {
	userID := c.Param("id")

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"message": "User not found"})
	}

	deposit := new(models.Transaction)
	if err := c.Bind(deposit); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Invalid deposit data"})
	}

	if deposit.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "Deposit amount must be greater than zero"})
	}

	user.AccountBalance += deposit.Amount

	if result := database.DB.Save(&user); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to update user balance"})
	}

	transaction := models.Transaction{
		UserID:    user.ID,
		Type:      "deposit",
		Amount:    deposit.Amount,
		Balance:   user.AccountBalance,
		CreatedAt: time.Now(),
	}

	if result := database.DB.Create(&transaction); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to log transaction"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Deposit successful",
		"balance": user.AccountBalance,
	})
}

func GetLoanHistory(c echo.Context) error {
	userID := c.Param("id")

	var transactions []models.Transaction
	if result := database.DB.Where("user_id = ?", userID).Find(&transactions); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"message": "Failed to fetch transactions"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":      "Transaction history fetched successfully",
		"transactions": transactions,
	})
}
