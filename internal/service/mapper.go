package service

import (
	"github.com/hunafazaky/event-booking-app/internal/dto"
	"github.com/hunafazaky/event-booking-app/internal/model"
)

func toUserResponse(user model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
