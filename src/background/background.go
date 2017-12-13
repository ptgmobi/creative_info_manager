package background

import (
	"fmt"
	"sync"
	"time"

	"cache"
	"creative_info"
	"db"
	"util"
)

func BatchRequestSize(cInfos []creative_info.CreativeInfo) []creative_info.CreativeInfo {
	var wg sync.WaitGroup
	cInfoChan := make(chan creative_info.CreativeInfo, len(cInfos))

	for _, cInfo := range cInfos {
		wg.Add(1)
		go func(info creative_info.CreativeInfo) {
			defer wg.Done()
			size, err := util.GetResourceSize(info.Url, 2500)
			if err != nil || size <= 0 {
				fmt.Println("BatchRequestSize error : ", err, ", url: ", info.Url, ", size: ", size)
				info.FailTimes++
			}
			info.Size = size
			cInfoChan <- info
		}(cInfo)
	}

	wg.Wait()
	close(cInfoChan)

	var newInfos []creative_info.CreativeInfo
	for newInfo := range cInfoChan {
		newInfos = append(newInfos, newInfo)
	}

	return newInfos
}

func Init() {
	go func() {
		for {
			func() {
				cInfos, err := db.GetCreativeInfoWithNoSize()
				if err != nil {
					fmt.Println("[background] db.GetCreativeInfoWithNoSize error: ", err)
					return
				}

				newInfos := BatchRequestSize(cInfos)
				if err := db.BatchUpdateSize(newInfos); err != nil {
					fmt.Println("[background] db.BatchUpdateSize error: ", err)
					return
				}
				if err := cache.BatchUpdateSize(newInfos); err != nil {
					fmt.Println("[background] cache.BatchUpdateSize error: ", err)
					return
				}
			}()

			fmt.Println("background running")
			time.Sleep(time.Second * 10)
		}
	}()
}
