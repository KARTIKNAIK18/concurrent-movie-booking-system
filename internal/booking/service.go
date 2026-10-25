package booking


import(
	"context"
)

type Service struct{
	store BookingStore
}



func NewService(store BookingStore) *Service{
	return &Service{
		store: store,
	}
}


//methods of Service
func(s *Service) BookTicket(b Booking)(Booking,error){
	return s.store.BookTicket(b)
}


func (s *Service) ListBookings(movieID string) []Booking {
	return s.store.ListBookings(movieID)
}

func (s *Service) ConfirmSeat(ctx context.Context, sessionID string, userID string) (Booking, error) {
	return s.store.Confirm(ctx, sessionID, userID)
}

func (s *Service) ReleaseSeat(ctx context.Context, sessionID string, userID string) error {
	return s.store.Release(ctx, sessionID, userID)
}