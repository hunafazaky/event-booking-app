package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/model"
	"github.com/hunafazaky/event-booking-app/internal/response"
	"github.com/hunafazaky/event-booking-app/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var input model.InputSignUp
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.SignUp(input)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "user created successfully", user)
}

func (h *UserHandler) SignIn(c *gin.Context) {
	var input model.InputSignIn
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.SignIn(input)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "user signed successfully", user)
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.FromError(c, err)
		return
	}

	user, err := h.service.GetAuthUser(userID)
	if err != nil {
		response.FromError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "data user retrieved", user)
}
