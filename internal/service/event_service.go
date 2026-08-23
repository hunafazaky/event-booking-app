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
	// List, GetByID, GetByUser, Update, Delete — same pattern as UserService,
	// you'll add these yourself once Create/Update are reviewed.
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
