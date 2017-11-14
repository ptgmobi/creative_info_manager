package cache

import (
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Conf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var cachePool *redis.Pool

func Init(cf *Conf) {
	if len(cf.Host) <= 0 || len(cf.Port) <= 0 {
		panic("no redis host or port")
	}

	if _, err := strconv.Atoi(cf.Port); err != nil {
		panic("redis port not number: " + cf.Port)
	}

	func(host, port string) {
		cachePool = &redis.Pool{
			MaxIdle:     256,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialTimeout("tcp", host+":"+port,
					50*time.Millisecond, 50*time.Millisecond, 50*time.Millisecond)
				if err != nil {
					return nil, err
				}
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < 10*time.Second {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
		}
	}(cf.Host, cf.Port)
}

func getConn() redis.Conn {
	return cachePool.Get()
}

func GetCreativeId(cUrl string) (int64, error) {
	c := getConn()
	id, err := redis.String(c.Do("HGet", "creative_info", cUrl))
	if err != nil {
		return 0, err
	}
	cId, _ := strconv.ParseInt(id, 10, 64)
	return cId, nil
}

func SetCreativeId(cUrl string, cId int64) error {
	c := getConn()
	id := strconv.FormatInt(cId, 10)
	_, err := c.Do("HSet", "creative_info", cUrl, id)
	return err
}
