package booking


type MemoryStore struct {
	bookings map[string]Booking
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		bookings: map[string]Booking{},
	}
}

func (s *MemoryStore) Book(b Booking)  error {
	//if seat is taken err or  give teh seat
	if _, exist := s.bookings[b.SeatID]; exist {
		return ErrSeatAlredyBooked
	}
	s.bookings[b.SeatID]=b
	return nil
}

func (s *MemoryStore) ListBookings(movieID string) []Booking {
	var result[] Booking

	for _, b := range s.bookings {
		if b.MovieID == movieID{
			result=append(result, b)
		}
	}

	return result
}