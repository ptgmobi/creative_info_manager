package search

import (
	"encoding/json"
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
}

func NewResp(errMsg, cId string) *Resp {
	return &Resp{
		ErrMsg:     errMsg,
		CreativeId: cId,
	}
}

func (resp *Resp) WriteTo(w http.ResponseWriter) (int, error) {
	b, _ := json.Marshal(&resp)
	return w.Write(b)
}

func (s *Service) HandleCreativeId(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := r.ParseForm(); err != nil {
		s.l.Println("ParseForm error: ", err)
		if n, err := NewResp("server error", "").WriteTo(w); err != nil {
			s.l.Println("[search] server error, resp write: ", n, ", error:", err)
		}
		return
	}

	cUrl, err := url.QueryUnescape(r.Form.Get("creative_url"))
	if err != nil || len(cUrl) == 0 {
		if n, err := NewResp("empty creative_url", "").WriteTo(w); err != nil {
			s.l.Println("[search] empty creative_url, resp write: ", n, ", error:", err)
		}
		return
	}

	if cId, err := cache.GetCreativeId(cUrl); err == nil && len(cId) > 0 {
		if n, err := NewResp("", cId).WriteTo(w); err != nil {
			s.l.Println("[search] cId in cache, cUrl: ", cUrl, ", resp write: ", n, ", error:", err)
		}
		return
	} else {
		if cId, err := db.GetCreativeId(cUrl); err == nil && len(cId) > 0 {
			if err := cache.SetCreativeId(cUrl, cId); err != nil {
				s.l.Println("[search] cache.SetCreativeId error, cUrl: ", cUrl, ", error: ", err)
			}
			if n, err := NewResp("", cId).WriteTo(w); err != nil {
				s.l.Println("[search] cId in db, cUrl: ", cUrl, ", resp write: ", n, ", error:", err)
			}
			return
		} else {
			s.l.Println("[search] db.GetCreativeId, cUrl: ", cUrl, ", err: ", err)
			if n, err := NewResp("database error", "").WriteTo(w); err != nil {
				s.l.Println("[search] database error, cUrl: ", cUrl, ", resp write: ", n, ", error:", err)
			}
			return
		}
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
