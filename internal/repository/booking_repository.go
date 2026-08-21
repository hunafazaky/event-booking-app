package repository

import (
	"github.com/hunafazaky/event-booking-app/internal/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *model.Booking) error
	Delete(booking *model.Booking) error
	FindByID(id string) (*model.Booking, error)
	FindByUserID(userID uint) (bookings []model.Booking, err error)
	ExistsByUserAndEvent(userID, eventID uint) (status bool, err error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) Delete(booking *model.Booking) error {
	return r.db.Delete(booking).Error
}

func (r *bookingRepository) FindByID(id string) (*model.Booking, error) {
	var booking model.Booking
	return &booking, r.db.First(&booking, id).Error
}

func (r *bookingRepository) FindByUserID(userID uint) (bookings []model.Booking, err error) {
	err = r.db.
		Preload("Event").
		Preload("Event.User", userSummary).
		Where("user_id = ?", userID).
		Find(&bookings).Error
	return
}

func (r *bookingRepository) ExistsByUserAndEvent(userID, eventID uint) (status bool, err error) {
	err = r.db.
		Where("user_id = ? AND event_id = ?", userID, eventID).
		First(&model.Booking{}).Error
	return err == nil, err
}
