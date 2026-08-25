package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hunafazaky/event-booking-app/internal/apperror"
	"github.com/hunafazaky/event-booking-app/internal/dto"
	"github.com/hunafazaky/event-booking-app/internal/model"
	"github.com/hunafazaky/event-booking-app/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	SignUp(input model.InputSignUp) (*dto.UserResponse, error)
	SignIn(input model.InputSignIn) (*dto.SignInResponse, error)
	GetAuthUser(userID uint) (*dto.UserResponse, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &userService{repo: repo, jwtSecret: jwtSecret}
}

func (s *userService) SignUp(input model.InputSignUp) (*dto.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperror.Internal("Failed to Hash Password", err)
	}

	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(&user); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("Email already registered")
		}
		return nil, apperror.Internal("Failed to create User", err)
	}

	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return &response, nil
}

func (s *userService) SignIn(input model.InputSignIn) (*dto.SignInResponse, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return nil, apperror.Unauthorized("Invalid Email or Password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, apperror.Unauthorized("Invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, apperror.Internal("Failed to sign token", err)
	}

	userResponse := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	response := dto.SignInResponse{
		Token: tokenString,
		User:  userResponse,
	}

	return &response, nil
}

func (s *userService) GetAuthUser(userID uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, mapLookupError(err, "user not found", "failed to load user")
	}

	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return &response, nil
}
