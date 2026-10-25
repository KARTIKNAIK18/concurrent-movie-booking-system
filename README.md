# Cinema Booking

A small cinema seat-booking application built with Go, Redis, and a browser-based frontend. Users can browse movies, view seat availability, temporarily hold a seat, confirm the booking, or release the hold.

## Features

- Movie list with seat layout information.
- Live seat availability polling in the frontend.
- Two-minute seat holds with automatic Redis expiry.
- Booking confirmation and manual hold release.
- Atomic Redis seat claims using `SET NX`, so only one concurrent request can hold a seat.
- A concurrency test that runs 100,000 attempts against the same seat.

## Requirements

- Go 1.26 or newer
- Redis running on `localhost:6379`

The application currently connects to Redis at `localhost:6379`.

## Run the application

Start Redis, then run the server from the repository root:

```bash
go run ./cmd
```

Open <http://localhost:8081> in a browser.

The server serves the files in `frontend/` and exposes the booking API from the same origin.

## API

### List movies

```http
GET /movies
```

Example response:

```json
[
	{"id":"Inception","title":"Inception","rows":5,"seats_per_row":8},
	{"id":"Dune","title":"Dune :Part Two","rows":4,"seats_per_row":6}
]
```

### List booked seats

```http
GET /movies/{movieID}/seats
```

### Hold a seat

```http
POST /movies/{movieID}/seats/{seatID}/hold
Content-Type: application/json

{"user_id":"user-123"}
```

The response contains a `session_id` and an `expires_at` timestamp. A hold is valid for two minutes unless it is confirmed or released.

### Confirm a held seat

```http
PUT /sessions/{sessionID}/confirm
Content-Type: application/json

{"user_id":"user-123"}
```

### Release a held seat

```http
DELETE /sessions/{sessionID}
Content-Type: application/json

{"user_id":"user-123"}
```

## Development

Run the Go test suite with:

```bash
go test ./...
```

The Redis concurrency test requires Redis to be available at `localhost:6379`. It verifies that exactly one of 100,000 simultaneous attempts can hold the same seat.

## Project structure

```text
cmd/                 HTTP server and movie catalog
frontend/            Browser UI, JavaScript, and CSS
internal/booking/     Booking domain, service, handlers, and stores
internal/adapters/    Redis client adapter
internal/test/        Service and concurrency tests
internal/utils/       Shared HTTP response helpers
```

## Booking lifecycle

1. A client requests a seat hold with a user ID.
2. Redis atomically claims the seat for two minutes.
3. The client confirms the session, which makes the seat persistent.
4. The client can release the session instead, or Redis releases an unconfirmed hold when its TTL expires.
