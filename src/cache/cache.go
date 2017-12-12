package cache

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Conf struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var cachePool *redis.Pool

func Init(cf *Conf) {
	if len(cf.Host) == 0 || len(cf.Port) == 0 {
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

func GetCreativeInfo(cUrl string) (string, int64, error) {
	c := cachePool.Get()
	defer c.Close()
	cInfo, err := redis.String(c.Do("Get", cUrl))
	if err != nil {
		return "", 0, err
	}
	info := strings.Split(cInfo, "_")
	if len(info) != 2 {
		return "", 0, errors.New("invalid info")
	} else {
		cSize, err := strconv.ParseInt(info[1], 10, 64)
		return info[0], cSize, err
	}
}

func SetCreativeInfo(cUrl, cId string, cSize int64) error {
	c := cachePool.Get()
	defer c.Close()
	_, err := c.Do("Set", cUrl, cId+"_"+strconv.FormatInt(cSize, 10))
	return err
}
