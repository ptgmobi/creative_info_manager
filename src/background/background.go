package background

import (
	"log"
	"sync"
	"time"

	"github.com/brg-liuwei/gotools"

	"cache"
	"creative_info"
	"db"
	"util"
)

type Conf struct {
	Interval        int    `json:"interval"`
	LogPath         string `json:"log_path"`
	LogRotateBackup int    `json:"log_rotate_backup"`
	LogRotateLines  int    `json:"log_rotate_lines"`
}

type Service struct {
	conf *Conf
	l    *gotools.RotateLogger
}

func NewService(conf *Conf) (*Service, error) {
	l, err := gotools.NewRotateLogger(conf.LogPath, "[background]", log.LUTC|log.LstdFlags, conf.LogRotateBackup)
	if err != nil {
		log.Println("[background] create log err: ", err)
		return nil, err
	}
	l.SetLineRotate(conf.LogRotateLines)

	srv := &Service{
		conf: conf,
		l:    l,
	}

	return srv, nil
}

func (s *Service) BatchRequestSize(cInfos []creative_info.CreativeInfo) []creative_info.CreativeInfo {
	var wg sync.WaitGroup
	cInfoChan := make(chan creative_info.CreativeInfo, len(cInfos))

	for _, cInfo := range cInfos {
		wg.Add(1)
		go func(info creative_info.CreativeInfo) {
			defer wg.Done()
			size, err := util.GetResourceSize(info.Url, 5000)
			if err != nil || size <= 0 {
				s.l.Println("BatchRequestSize error : ", err, ", url: ", info.Url, ", size: ", size)
				size = 0
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

func (s *Service) LoopUpdateSize() {
	cInfos, err := db.GetCreativeInfoWithNoSize()
	if err != nil {
		s.l.Println("[background] GetCreativeInfoWithNoSize error: ", err)
		return
	}
	if len(cInfos) == 0 {
		s.l.Println("[background] GetCreativeInfoWithNoSize all urls has size, except for those fail more than 5 times")
		return
	}

	newInfos := s.BatchRequestSize(cInfos)
	if err := db.BatchUpdateSize(newInfos); err != nil {
		s.l.Println("[background] db.BatchUpdateSize error: ", err)
		return
	}
	if err := cache.BatchUpdateSize(newInfos); err != nil {
		s.l.Println("[background] cache.BatchUpdateSize error: ", err)
		return
	}
}

func (s *Service) Serve() {
	for {
		s.LoopUpdateSize()
		time.Sleep(time.Second * time.Duration(s.conf.Interval))
	}
}
