package service

import (
	"github.com/hunafazaky/event-booking-app/internal/dto"
	"github.com/hunafazaky/event-booking-app/internal/repository"
)

type EventService interface {
  Create(userID uint, input CreateEventInput) (*dto.EventResponse, error)
	Update(userID, eventID uint, input UpdateEventInput) (*dto.EventResponse, error)
	Delete(userID, eventID uint) error
	GetByID(id uint) (*dto.EventDetailResponse, error)
	GetByUser(userID uint) ([]dto.EventResponse, error)
	List(search string, page, limit int) ([]dto.EventResponse, dto.EventListMeta, error)


type eventService struct {
  repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
  return &eventService{repo: repo}
}
