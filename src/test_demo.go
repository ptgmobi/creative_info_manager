package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func getRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func getCreativeId(cUrl string, wg sync.WaitGroup) {
	defer wg.Done()
	cUrl += getRandomString(10)
	uri := "http://localhost:12121/get_creative_id?creative_url=" + cUrl

	resp, err := http.Get(uri)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println("getCreativeId http get error: ", err)
		return
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("getCreativeId read body error: ", err)
		return
	}
	log.Println("getCreativeId success: ", string(bytes))

}

func main() {
	baseUrl := "https://www.baidu.com/"
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go getCreativeId(baseUrl, wg)
	}
	wg.Wait()
}
