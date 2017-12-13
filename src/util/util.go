package util

import (
	"net/http"
	"time"
)

func GetResourceSize(url string, timeout int) (int64, error) {
	if timeout == 0 {
		timeout = 200
	}
	c := &http.Client{
		Timeout: time.Millisecond * time.Duration(timeout),
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
