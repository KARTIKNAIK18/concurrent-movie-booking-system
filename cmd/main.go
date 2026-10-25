package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/KARTIKNAIK18/concurrent-movie-booking-system/internal/adapters/redis"
	"github.com/KARTIKNAIK18/concurrent-movie-booking-system/internal/booking"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /movies", listMovies)

	mux.Handle("GET /", http.FileServer(http.Dir("frontend")))

	store := booking.NewRedditStore(redis.Client("localhost:6379"))
	svc := booking.NewService(store)

	bookingHandler := booking.NewHandler(svc)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeat)

	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler.ReleaseSession)

	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}

var movies = []movieResponse{
	{ID: "Inception", Title: "Inception", Rows: 5, SeatPerRow: 8},
	{ID: "Dune", Title: "Dune :Part Two", Rows: 4, SeatPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, http.StatusOK, movies)
}

func bookingHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "booking handler not implemented", http.StatusNotImplemented)
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

type movieResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Rows       int    `json:"rows"`
	SeatPerRow int    `json:"seats_per_row"`
}
