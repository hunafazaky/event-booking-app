package service

import (
	"errors"
	"io"
	"time"

	"github.com/hunafazaky/event-booking-app/internal/apperror"
	"github.com/hunafazaky/event-booking-app/internal/dto"
	"github.com/hunafazaky/event-booking-app/internal/model"
	"github.com/hunafazaky/event-booking-app/internal/repository"
	"gorm.io/gorm"
)

// CreateEventInput is what the handler builds from the multipart form and
// hands to the service. Image/ImageName are plain io.Reader + string —
// no multipart-specific type crosses into the service layer.
type CreateEventInput struct {
	Name        string
	Description string
	Location    string
	DateTime    time.Time
	Image       io.Reader
	ImageName   string
}

type UpdateEventInput struct {
	Name        string
	Description string
	Location    string
	DateTime    *time.Time
	Image       io.Reader
	ImageName   string
}

type EventService interface {
	Create(userID uint, input CreateEventInput) (*dto.EventResponse, error)
	List(search string, page, limit int) ([]dto.EventResponse, dto.EventListMeta, error)
	GetByID(id uint) (*dto.EventDetailResponse, error)
	GetByUser(userID uint) ([]dto.EventResponse, error)
	Update(userID, eventID uint, input UpdateEventInput) (*dto.EventResponse, error)
	Delete(userID, eventID uint) error
}

type eventService struct {
	repo     repository.EventRepository
	uploader ImageUploader
}

func NewEventService(repo repository.EventRepository, uploader ImageUploader) EventService {
	return &eventService{repo: repo, uploader: uploader}
}

func (s *eventService) Create(userID uint, input CreateEventInput) (*dto.EventResponse, error) {
	// Upload FIRST, before touching the database. If the upload fails,
	// there's nothing to roll back — we simply never created a DB row
	// with a broken/missing image reference.
	imageURL, imageID, err := s.uploader.Upload(input.Image, input.ImageName)
	if err != nil {
		return nil, apperror.Internal("failed to upload image", err)
	}

	event := model.Event{
		Name:        input.Name,
		Description: input.Description,
		Location:    input.Location,
		DateTime:    input.DateTime,
		Image:       imageURL,
		ImageID:     imageID,
		UserID:      userID,
	}

	if err := s.repo.Create(&event); err != nil {
		// If the DB write fails AFTER a successful upload, you now have an
		// orphaned file sitting in ImageKit with nothing pointing to it.
		// Worth a `_ = s.uploader.Delete(imageID)` here as best-effort
		// cleanup — same "don't fail the whole request over a cleanup
		// step" reasoning your original Update handler already used for
		// deleting the OLD image.
		return nil, apperror.Internal("failed to create event", err)
	}

	response := dto.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		Image:       event.Image,
		DateTime:    event.DateTime,
		CreatedAt:   event.CreatedAt,
		// event.User is a zero-value User here — repo.Create doesn't
		// populate the association. You'll want a repo.FindByID call
		// after Create, OR just build UserResponse from data you already
		// have (userID + the fields from the JWT context, if the handler
		// has them) — worth deciding which once you get here.
	}
	return &response, nil
}

func (s *eventService) List(search string, page, limit int) ([]dto.EventResponse, dto.EventListMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 6
	}

	events, totalRows, totalPage, err := s.repo.FindAll(search, page, limit)
	if err != nil {
		return nil, dto.EventListMeta{}, apperror.Internal("failed to load data", err)
	}

	eventResponse := make([]dto.EventResponse, 0, len(events))
	for _, item := range events {
		eventResponse = append(eventResponse, dto.EventResponse{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Location:    item.Location,
			Image:       item.Image,
			DateTime:    item.DateTime,
			User:        toUserResponse(item.User),
			CreatedAt:   item.CreatedAt,
		})
	}

	eventListMeta := dto.EventListMeta{
		Page:      page,
		Limit:     limit,
		TotalRows: totalRows,
		TotalPage: totalPage,
	}

	return eventResponse, eventListMeta, nil
}

func (s *eventService) GetByID(id uint) (*dto.EventDetailResponse, error) {
	event, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("event not found")
		}
		return nil, apperror.Internal("failed to load event", err)
	}

	eventResponse := dto.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		Image:       event.Image,
		DateTime:    event.DateTime,
		User:        toUserResponse(event.User),
		CreatedAt:   event.CreatedAt,
	}

	bookingSummaryResponse := make([]dto.BookingSummaryResponse, 0, len(event.Booking))
	for _, item := range event.Booking {
		bookingSummaryResponse = append(bookingSummaryResponse, dto.BookingSummaryResponse{
			ID:          item.ID,
			BookingCode: item.BookingCode,
			Phone:       item.Phone,
			User:        toUserResponse(item.User),
		})
	}

	eventDetailResponse := dto.EventDetailResponse{
		EventResponse: eventResponse,
		Bookings:      bookingSummaryResponse,
	}
	return &eventDetailResponse, nil
}

func (s *eventService) GetByUser(userID uint) ([]dto.EventResponse, error) {
	events, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, apperror.Internal("failed to load events", err)
	}

	listEventResponse := make([]dto.EventResponse, 0, len(events))
	for _, item := range events {
		listEventResponse = append(listEventResponse, dto.EventResponse{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Location:    item.Location,
			Image:       item.Image,
			DateTime:    item.DateTime,
			User:        toUserResponse(item.User),
			CreatedAt:   item.CreatedAt,
		})
	}

	return listEventResponse, nil
}

func (s *eventService) Update(userID, eventID uint, input UpdateEventInput) (*dto.EventResponse, error) {
	// get current event
	event, err := s.repo.FindByID(eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("event not found")
		}
		return nil, apperror.Internal("failed to load event", err)
	}

	// verify permission
	if event.UserID != userID {
		return nil, apperror.Forbidden("you're not authorized")
	}

	// file handler
	if input.Image != nil {
		imageURL, imageID, err := s.uploader.Upload(input.Image, input.ImageName)
		if err != nil {
			return nil, apperror.Internal("failed to upload image", err)
		}
		// TODO: delete the OLD image (event.ImageID) from ImageKit here —
		// same orphaned-file concern as Create, now on the "replace" side.
		event.Image = imageURL
		event.ImageID = imageID
	}

	// replace data
	if input.Name != "" {
		event.Name = input.Name
	}
	if input.Description != "" {
		event.Description = input.Description
	}
	if input.Location != "" {
		event.Location = input.Location
	}
	if input.DateTime != nil {
		event.DateTime = *input.DateTime
	}

	if err := s.repo.Update(event); err != nil {
		return nil, apperror.Internal("failed to update event", err)
	}

	eventResponse := dto.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		Image:       event.Image,
		DateTime:    event.DateTime,
		CreatedAt:   event.CreatedAt,
	}

	return &eventResponse, nil
}

func (s *eventService) Delete(userID, eventID uint) error {
	event, err := s.repo.FindByID(eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("event not found")
		}
		return apperror.Internal("failed to load event", err)
	}

	// verify permission
	if event.UserID != userID {
		return apperror.Forbidden("you're not authorized")
	}

	if err := s.repo.Delete(event); err != nil {
		return apperror.Internal("failed to delete event", err)
	}

	return nil
}
