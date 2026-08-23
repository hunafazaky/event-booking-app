package service

import (
	"io"
	"time"

	"github.com/hunafazaky/event-booking-app/internal/apperror"
	"github.com/hunafazaky/event-booking-app/internal/dto"
	"github.com/hunafazaky/event-booking-app/internal/model"
	"github.com/hunafazaky/event-booking-app/internal/repository"
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

type EventService interface {
	Create(userID uint, input CreateEventInput) (*dto.EventResponse, error)
	List(search string, page, limit int) ([]dto.EventResponse, dto.EventListMeta, error)
	GetByID(id uint) (*dto.EventDetailResponse, error)
	// GetByUser(userID uint) ([]dto.EventResponse, error)
	// Update(userID, eventID uint, input UpdateEventInput) (*dto.EventResponse, error)
	// Delete(userID, eventID uint) error
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
	url, fileID, err := s.uploader.Upload(input.Image, input.ImageName)
	if err != nil {
		return nil, apperror.Internal("failed to upload image", err)
	}

	event := model.Event{
		Name:        input.Name,
		Description: input.Description,
		Location:    input.Location,
		DateTime:    input.DateTime,
		Image:       url,
		ImageID:     fileID,
		UserID:      userID,
	}

	if err := s.repo.Create(&event); err != nil {
		// If the DB write fails AFTER a successful upload, you now have an
		// orphaned file sitting in ImageKit with nothing pointing to it.
		// Worth a `_ = s.uploader.Delete(fileID)` here as best-effort
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
		return nil, dto.EventListMeta{}, apperror.Internal("Failed to load data", err)
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
			User: dto.UserResponse{
				ID:    item.User.ID,
				Name:  item.User.Name,
				Email: item.User.Email,
			},
			CreatedAt: item.CreatedAt,
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
		return nil, apperror.Internal("Failed to load data", err)
	}

	eventResponse := dto.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		Image:       event.Image,
		DateTime:    event.DateTime,
		User: dto.UserResponse{
			ID:    event.User.ID,
			Name:  event.User.Name,
			Email: event.User.Email,
		},
		CreatedAt: event.CreatedAt,
	}

	bookingSummaryResponse := make([]dto.BookingSummaryResponse, 0, len(event.Booking))
	for _, item := range event.Booking {
		bookingSummaryResponse = append(bookingSummaryResponse, dto.BookingSummaryResponse{
			ID:          item.ID,
			BookingCode: item.BookingCode,
			Phone:       item.Phone,
			User: dto.UserResponse{
				ID:    event.User.ID,
				Name:  event.User.Name,
				Email: event.User.Email,
			},
		})
	}

	eventDetailResponse := dto.EventDetailResponse{
		EventResponse: eventResponse,
		Bookings:      bookingSummaryResponse,
	}
	return &eventDetailResponse, nil
}
