package search

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/brg-liuwei/gotools"

	"cache"
	"db"
)

type Conf struct {
	LogPath         string `json:"log_path"`
	LogRotateBackup int    `json:"log_rotate_backup"`
	LogRotateLines  int    `json:"log_rotate_lines"`
}

type Service struct {
	conf *Conf
	l    *gotools.RotateLogger
}

type Resp struct {
	ErrMsg     string `json:"err_msg"`
	CreativeId string `json:"creative_id"`
	Size       int64  `json:"size"`
}

func NewResp(errMsg, cId, cType string, cSize int64) *Resp {
	if len(cId) > 0 {
		switch cType {
		case "1":
			cId = "img." + cId
		case "2":
			cId = "mp4." + cId // 我们视频素材暂时只有mp4
		default:
			log.Println("unknown creative type")
		}
	}

	return &Resp{
		ErrMsg:     errMsg,
		CreativeId: cId,
		Size:       cSize,
	}
}

func (resp *Resp) WriteTo(w http.ResponseWriter) (int, error) {
	b, _ := json.Marshal(&resp)
	return w.Write(b)
}

func GetInfoFromDbAndSetCache(cUrl, cType string) (string, int64, error) {
	cId, cSize, err := db.GetCreativeInfo(cUrl, cType)
	if err == nil || len(cId) == 0 {
		return "", 0, fmt.Errorf("db.GetCreativeInfo error: %v", err)
	}

	err = cache.SetCreativeInfo(cUrl, cId, cSize)
	if err != nil {
		log.Println("GetInfoFromDbAndSetCache SetCreativeInfo error: ", err, ", cUrl: ", cUrl, ", size: ", cSize)
	}

	return cId, cSize, nil

}

func (s *Service) HandleCreativeId(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := r.ParseForm(); err != nil {
		s.l.Println("ParseForm error: ", err)
		if n, err := NewResp("server error", "", "", 0).WriteTo(w); err != nil {
			s.l.Println("[search] server error, resp write: ", n, ", error:", err)
		}
		return
	}

	cUrl, err := url.QueryUnescape(r.Form.Get("creative_url"))
	if err != nil || len(cUrl) == 0 {
		s.l.Println("[search] can't get creative_url, err: ", err)
		if n, err := NewResp("can't get creative_url", "", "", 0).WriteTo(w); err != nil {
			s.l.Println("[search] can't get creative_url, resp write: ", n, ", error:", err)
		}
		return
	}

	cType := r.Form.Get("type")
	if len(cType) == 0 {
		cType = "1"
	}

	cId, cSize, err := cache.GetCreativeInfo(cUrl)
	if err == nil && len(cId) > 0 {
		if n, err := NewResp("", cId, cType, cSize).WriteTo(w); err != nil {
			s.l.Println("[search] fail to response cache cId, cUrl: ", cUrl, ", resp write: ", n, ", error:", err)
		}
		return
	}

	cId, cSize, err = GetInfoFromDbAndSetCache(cUrl, cType)
	if err != nil || len(cId) == 0 {
		s.l.Println("[search] GetInfoFromDbAndSetCache error: ", err, ",  or empty cId: ", cId, ", cUrl: ", cUrl)
		if n, err := NewResp("server error", "", "", 0).WriteTo(w); err != nil {
			s.l.Println("[search] fail to response server error, cUrl: ", cUrl, ", resp write: ", n, ", error:", err)
		}
		return
	}

	if n, err := NewResp("", cId, cType, cSize).WriteTo(w); err != nil {
		s.l.Println("[search] fail to response db cId, cUrl: ", cUrl, ", resp write: ", n, ", error:", err)
	}
}

func NewService(conf *Conf) (*Service, error) {
	l, err := gotools.NewRotateLogger(conf.LogPath, "[search]", log.LUTC|log.LstdFlags, conf.LogRotateBackup)
	if err != nil {
		log.Println("[Search] create log err: ", err)
		return nil, err
	}
	l.SetLineRotate(conf.LogRotateLines)

	srv := &Service{
		conf: conf,
		l:    l,
	}

	return srv, nil
}

func (s *Service) Serve() {
	http.HandleFunc("/get_creative_id", s.HandleCreativeId)
	panic(http.ListenAndServe(":12121", nil))
}
