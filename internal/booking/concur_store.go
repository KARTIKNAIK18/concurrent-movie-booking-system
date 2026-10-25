package booking

import "sync"

type ConcurrentStore struct {
	bookings map[string]Booking
	mu sync.RWMutex
}

func NewConurrentStore() *ConcurrentStore {
	return &ConcurrentStore{
		bookings: map[string]Booking{},
	}
} 

func (s *ConcurrentStore) Book(b Booking) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	//if seat is taken err or  give teh seat
	if _, exist := s.bookings[b.SeatID]; exist {
		return ErrSeatAlredyBooked
	}
	s.bookings[b.SeatID] = b
	return nil
}

func (s *ConcurrentStore) ListBookings(movieID string) []Booking {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []Booking

	for _, b := range s.bookings {
		if b.MovieID == movieID {
			result = append(result, b)
		}
	}

	return result
}