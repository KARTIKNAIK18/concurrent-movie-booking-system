package booking

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	uuid "github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
)


const defaultHoldTTL= 2*time.Minute


type Redistore struct{
	rdb *redis.Client
}


func NewRedditStore(rdb *redis.Client) *Redistore{
	return &Redistore{rdb: rdb}
}

func (s *Redistore) BookTicket(b Booking) (Booking,error) {
	res,err:=s.Booknow(b)
	
	if err != nil {
		return Booking{}, err
	}

	log.Printf("Session Booked %v", res)
	return res,nil
}



func sessionKey(id string) string {
	return fmt.Sprintf("session:%s", id)
}



func (s *Redistore)Booknow(b Booking)(Booking,error){
	Id:=uuid.New().String()
	ttl:=time.Now()
	ctx:=context.Background()
	key := fmt.Sprintf("seat:%s:%s", b.MovieID, b.SeatID)

	
	
	b.ID=Id
	_,val:=json.Marshal(b)
	

	res:=s.rdb.SetArgs(ctx,key,val,redis.SetArgs{
		Mode: "NX",
		TTL: defaultHoldTTL,
	})
	

	if confirm:=res.Val(); confirm!="OK"{
		return Booking{}, ErrSeatAlredyBooked
	}

	s.rdb.Set(ctx,sessionKey(Id),key,defaultHoldTTL)

	return Booking{
		ID: Id,
		MovieID: b.MovieID,
		SeatID: b.SeatID,
		UserID: b.UserID,
		Status: "held",
		ExpiresAt: ttl.Add(defaultHoldTTL),
	},nil
}


func (s *Redistore) ListBookings(MovieID string) [] Booking{
	pattern := fmt.Sprintf("seat:%s:*", MovieID)
	Session:=[] Booking{}

	ctx:=context.Background()
	data:=s.rdb.Scan(ctx,0,pattern,0).Iterator()

	for data.Next(ctx) {
		val,Err:= s.rdb.Get(ctx,data.Val()).Result()

		if Err != nil {
			continue
		}

		data,err:=ParseSession(val)

		if err != nil {
			continue
		}

		Session=append(Session, data)

	}

	return Session
}



func ParseSession(val string) (Booking,error){

	var data Booking

	if err:=json.Unmarshal([]byte(val), &data); err!=nil{
		return Booking{},err
	}

	return Booking{
		ID:      data.ID,
		MovieID: data.MovieID,
		SeatID:  data.SeatID,
		UserID:  data.UserID,
		Status:  data.Status,
	}, nil

}

func (s *Redistore) Confirm(ctx context.Context, sessionID string, userID string) (Booking, error) {
	session, sk, err := s.getSession(ctx, sessionID, userID)
	if err != nil {
		return Booking{}, err
	}

	s.rdb.Persist(ctx, sk)
	s.rdb.Persist(ctx, sessionKey(sessionID))

	session.Status = "confirmed"
	data := Booking{
		ID:      session.ID,
		MovieID: string(session.MovieID),
		SeatID:  session.SeatID,
		UserID:  session.UserID,
		Status:  "confirmed",
	}
	val, _ := json.Marshal(data)
	s.rdb.Set(ctx, sk, val, 0)

	return session, nil
}


func (s *Redistore) getSession(ctx context.Context, sessionID string, userID string) (Booking, string, error) {
	sk, err := s.rdb.Get(ctx, sessionKey(sessionID)).Result()
	if err != nil {
		return Booking{}, "", err
	}

	val, err := s.rdb.Get(ctx, sk).Result()
	if err != nil {
		return Booking{}, "", err
	}

	session, err := ParseSession(val)
	if err != nil {
		return Booking{}, "", err
	}

	return session, sk, nil
}


func (s *Redistore) Release(ctx context.Context, sessionID string, userID string) error {
	_, sk, err := s.getSession(ctx, sessionID, userID)
	if err != nil {
		return err
	}

	s.rdb.Del(ctx, sk, sessionKey(sessionID))
	return nil
}



