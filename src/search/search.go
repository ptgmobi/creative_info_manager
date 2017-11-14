package search

import (
	"encoding/json"
	"log"
	"net/http"

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
	Err string `json:"error"`
	CId int64  `json:"creative_id"`
}

func NewResp(errMsg string, cId int64) *Resp {
	return &Resp{
		Err: errMsg,
		CId: cId,
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
	w.Header().Set("Content-type", "application/json; charset=utf-8")
	r.ParseForm()

	cUrl := r.Form.Get("creative_url")
	if len(cUrl) <= 0 {
		if n, err := NewResp("empty creative_url", 0).WriteTo(w); err != nil {
			s.l.Println("[search] empty creative_url resp write: ", n, ", error:", err)
			return
		}
	}

	if cId, err := cache.GetCreativeId(cUrl); err == nil && cId > 0 {
		if n, err := NewResp("", cId).WriteTo(w); err != nil {
			s.l.Println("[search] cId in cache resp write: ", n, ", error:", err)
			return
		}
	} else {
		if cId, err := db.GetCreativeId(cUrl); err == nil && cId > 0 {
			if err := cache.SetCreativeId(cUrl, cId); err != nil {
				s.l.Println("[search] cache.SetCreativeId err: ", err)
			}
			if n, err := NewResp("", cId).WriteTo(w); err != nil {
				s.l.Println("[search] cId in db resp write: ", n, ", error:", err)
				return
			}
		} else {
			s.l.Println("[search] db.GetCreativeId err: ", err)
			if n, err := NewResp("database error", 0).WriteTo(w); err != nil {
				s.l.Println("[search] database error resp write: ", n, ", error:", err)
				return
			}
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
