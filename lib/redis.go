package lib

import "github.com/garyburd/redigo/redis"

func NewRedisPoll(redisURI string, db int) *redis.Pool {
	return newPool(redisURI, db)
}
func newPool(url string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", url, redis.DialDatabase(db))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
