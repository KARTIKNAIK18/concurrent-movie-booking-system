package booking

import (
	"errors"
	"time"
	"context"
)

var (
	ErrSeatAlredyBooked = errors.New("Seat is Already Taken")
)

// Booking represents a confirmed seat reservation.
type Booking struct {
	ID      string
	MovieID string
	SeatID  string
	UserID  string
	Status  string
	ExpiresAt time.Time
}

type BookingStore interface {
	BookTicket(b Booking) (Booking,error)
	ListBookings(movie_id string) [] Booking

	Confirm(ctx context.Context, sessionID string, userID string) (Booking, error)
	Release(ctx context.Context, sessionID string, userID string) error
}