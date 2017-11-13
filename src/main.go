package main

import (
	"github.com/brg-liuwei/gotools"

	"cache"
	"db"
	"search"
)

type Conf struct {
	DbConf     db.Conf     `json:"mysql_config"`
	CacheConf  cache.Conf  `json:"redis_config"`
	SearchConf search.Conf `json:"search_config"`
}

var conf Conf

func startSearchService(cf *search.Conf) {
	searchService, err := search.NewService(cf)
	if err != nil {
		panic(err)
	}
	searchService.Serve()
}

func main() {
	if err := gotools.DecodeJsonFile("conf/creative.conf", &conf); err != nil {
		panic(err)
	}

	cache.Init(&conf.CacheConf)
	db.Init(&conf.DbConf)

	startSearchService(&conf.SearchConf)
}
