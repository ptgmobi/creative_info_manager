package main

import (
	"github.com/brg-liuwei/gotools"

	"background"
	"cache"
	"db"
	"search"
)

type Conf struct {
	DbConf     db.Conf         `json:"mysql_config"`
	CacheConf  cache.Conf      `json:"redis_config"`
	SearchConf search.Conf     `json:"search_config"`
	BgConf     background.Conf `json:"background_config"`
}

var conf Conf

func startSearchService(cf *search.Conf) {
	searchService, err := search.NewService(cf)
	if err != nil {
		panic(err)
	}
	searchService.Serve()
}

func startBackgroundService(cf *background.Conf) {
	bgService, err := background.NewService(cf)
	if err != nil {
		panic(err)
	}
	bgService.Serve()
}

func main() {
	if err := gotools.DecodeJsonFile("conf/creative.conf", &conf); err != nil {
		panic(err)
	}

	cache.Init(&conf.CacheConf)
	db.Init(&conf.DbConf)

	startBackgroundService(&conf.BgConf)
	startSearchService(&conf.SearchConf)
}
