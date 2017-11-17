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

	cachePool = &redis.Pool{
		MaxIdle:     256,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", cf.Host+":"+cf.Port,
				100*time.Millisecond, 100*time.Millisecond, 100*time.Millisecond)
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
}

func GetCreativeId(cUrl string) (string, error) {
	c := cachePool.Get()
	defer c.Close()
	cId, err := redis.String(c.Do("HGet", "creative_info", cUrl))
	if err != nil {
		return "", err
	}
	return cId, nil
}

func SetCreativeId(cUrl, cId string) error {
	c := cachePool.Get()
	defer c.Close()
	_, err := c.Do("HSet", "creative_info", cUrl, cId)
	return err
}
