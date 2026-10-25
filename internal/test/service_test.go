package test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/google/uuid"

	"github.com/KARTIKNAIK18/concurrent-movie-booking-system/internal/adapters/redis"
	"github.com/KARTIKNAIK18/concurrent-movie-booking-system/internal/booking"
)

func TestConcurrentBooking_ExactlyOneWins(t *testing.T) {
	// store := booking.NewConurrentStore()
	// svc := booking.NewService(store)
	host := "localhost:6379"
	store := booking.NewRedditStore(redis.Client(host))
	svc := booking.NewService(store)

	const numGoroutines = 100_000 // 100k users trying to book a seat at the same time

	var (
		successes atomic.Int64
		failures  atomic.Int64
		wg        sync.WaitGroup
	)

	wg.Add(numGoroutines)
	for i := range numGoroutines {
		go func(userNum int) {
			defer wg.Done()
			_,err := svc.BookTicket(booking.Booking{
				MovieID: "screen-1",
				SeatID:  "A1",
				UserID:  uuid.New().String(),
			})
			if err == nil {
				successes.Add(1)
			} else {
				failures.Add(1)
			}
		}(i)
	}
	wg.Wait()

	if got := successes.Load(); got != 1 {
		t.Errorf("expected exactly 1 success, got %d", got)
	}
	if got := failures.Load(); got != int64(numGoroutines-1) {
		t.Errorf("expected %d failures, got %d", numGoroutines-1, got)
	}
}
