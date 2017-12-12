package util

import (
	"net/http"
	"time"
)

func GetResourceSize(url string) (int64, error) {
	c := &http.Client{
		Timeout: time.Millisecond * 1000,
	}
	resp, err := c.Get(url)
	if resp != nil {
		resp.Body.Close()
	}
	if err != nil {
		return 0, err
	}
	return resp.ContentLength, nil
}
