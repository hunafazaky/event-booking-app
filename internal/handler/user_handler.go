package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hunafazaky/event-booking-app/internal/response"
	"github.com/hunafazaky/event-booking-app/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// SignUp godoc
// @Summary      Register a new user
// @Description  Creates a new user account.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body      service.SignUpInput  true  "Sign up payload"
// @Success      201    {object}  response.Envelope{data=dto.UserResponse}
// @Failure      400    {object}  response.Envelope
// @Failure      409    {object}  response.Envelope
// @Router       /auth/signup [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	var input service.SignUpInput
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
	var input service.SignInInput
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
