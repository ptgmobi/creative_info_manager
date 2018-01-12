package dump

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/brg-liuwei/gotools"

	"cache"
	"creative_info"
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
	ErrMsg        string                       `json:"err_msg"`
	CreativeInfos []creative_info.CreativeInfo `json:"creative_info"`
}

func NewResp(errMsg string, cInfos []creative_info.CreativeInfo) *Resp {
	return &Resp{
		ErrMsg:        errMsg,
		CreativeInfos: cInfos,
	}
}

func (resp *Resp) WriteTo(w http.ResponseWriter) (int, error) {
	b, _ := json.Marshal(&resp)
	return w.Write(b)
}

func (s *Service) HandleDump(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := r.ParseForm(); err != nil {
		s.l.Println("ParseForm error: ", err)
		if n, err := NewResp("server error", nil).WriteTo(w); err != nil {
			s.l.Println("[dump] server error, resp write: ", n, ", error:", err)
		}
		return
	}

	var cInfos []creative_info.CreativeInfo
	var err error
	cIds := r.Form.Get("id")
	if len(cIds) == 0 {
		cInfos, err = cache.GetScanInfos(rand.Intn(db.GetMaxId()))
	} else {
		cIds = strings.Replace(cIds, "img.", "", -1)
		cIds = strings.Replace(cIds, "mp4.", "", -1)
		cInfos, err = db.GetCreativeInfoByIds(cIds)
	}
	if err != nil {
		s.l.Println("[dump] Get data from db error, cIds: ", cIds, ", err: ", err)
		if n, err := NewResp("database error", nil).WriteTo(w); err != nil {
			s.l.Println("[dump] database error, cIds: ", cIds, ", resp write: ", n, ", error:", err)
		}
		return
	}

	if n, err := NewResp("", cInfos).WriteTo(w); err != nil {
		s.l.Println("[dump] creative info resp write error, cIds: ", cIds, ", resp write: ", n, ", error:", err)
	}

	return
}

func NewService(conf *Conf) (*Service, error) {
	l, err := gotools.NewRotateLogger(conf.LogPath, "[dump]", log.LUTC|log.LstdFlags, conf.LogRotateBackup)
	if err != nil {
		log.Println("[Dump] create log err: ", err)
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
	http.HandleFunc("/dump", s.HandleDump)
	panic(http.ListenAndServe(":12345", nil))
}
