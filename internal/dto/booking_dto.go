package dto

// BookingResponse is the full shape of a booking, returned by
// POST /bookings and GET /bookings — includes which event it's for.
type BookingResponse struct {
	ID          uint          `json:"id"`
	BookingCode string        `json:"booking_code"`
	Phone       string        `json:"phone"`
	Event       EventResponse `json:"event"`
}

// BookingSummaryResponse is the reduced shape of a booking as shown
// nested inside EventDetailResponse. It shows who booked, not the full
// event again (that would be circular — the event is the parent here).
type BookingSummaryResponse struct {
	ID          uint         `json:"id"`
	BookingCode string       `json:"booking_code"`
	Phone       string       `json:"phone"`
	User        UserResponse `json:"user"`
}
