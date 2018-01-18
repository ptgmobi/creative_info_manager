package cache

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"

	"creative_info"
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
	value := cId + "_" + strconv.FormatInt(cSize, 10)
	_, err := c.Do("Set", cUrl, value)
	return err
}

func BatchUpdateSize(infos []creative_info.CreativeInfo) error {
	c := cachePool.Get()
	defer c.Close()

	for _, info := range infos {
		if info.Size > 0 {
			value := info.Id + "_" + strconv.FormatInt(info.Size, 10)
			if err := c.Send("Set", info.Url, value); err != nil {
				return errors.New("Set (" + info.Url + " " + value + ") error: " + err.Error())
			}
		}
	}

	if err := c.Flush(); err != nil {
		return errors.New("FLUSH creative info to redis error: " + err.Error())
	}

	return nil
}

func GetScanInfos(cursor int) ([]creative_info.CreativeInfo, error) {
	c := cachePool.Get()
	defer c.Close()

	arr, err := redis.MultiBulk(c.Do("SCAN", cursor))
	if err != nil {
		return nil, err
	}
	if len(arr) < 2 {
		return nil, errors.New("No url to show.")
	} else {
		keys, err := redis.Strings(arr[1], nil)
		if err != nil {
			return nil, err
		}
		var cInfos []creative_info.CreativeInfo
		for _, url := range keys {
			cInfos = append(cInfos, creative_info.CreativeInfo{Url: url})
		}
		return cInfos, nil
	}
}

func GetRandomKey() ([]creative_info.CreativeInfo, error) {
	c := cachePool.Get()
	defer c.Close()

	url, err := redis.String(c.Do("RANDOMKEY"))
	if err != nil || len(url) == 0 {
		return nil, fmt.Errorf("GetRandomKey error: %v", err)
	}
	var cInfos []creative_info.CreativeInfo
	cInfos = append(cInfos, creative_info.CreativeInfo{Url: url})
	return cInfos, nil
}
