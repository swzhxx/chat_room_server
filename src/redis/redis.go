package db

import "github.com/garyburd/redigo/redis"

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     8,
		MaxActive:   0,
		IdleTimeout: 100,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}
}

func GetRedisConn() (conn redis.Conn) {
	return pool.Get()
}
