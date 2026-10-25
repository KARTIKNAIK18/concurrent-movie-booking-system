package redis

import (
	"context"
	"log"

	redisCli "github.com/redis/go-redis/v9"
)



func Client(addrs string) *redisCli.Client{
	rdb := redisCli.NewClient(&redisCli.Options{Addr:addrs,
	})

	con:=context.Background()

	if err:=rdb.Ping(con).Err(); err!=nil{
		log.Fatalf("redis ping %v", err)
	}

	return rdb
}






