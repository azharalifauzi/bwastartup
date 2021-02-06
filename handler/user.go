package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterInput

	err := c.ShouldBindJSON(&input)

	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Register Account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	createdUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.ApiResponse("Register Account failed", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formattedUser := user.FormatUser(createdUser, "tokencobacoba")

	response := helper.ApiResponse("Account has been registered", http.StatusCreated, "success", formattedUser)

	c.JSON(http.StatusCreated, response)
}
