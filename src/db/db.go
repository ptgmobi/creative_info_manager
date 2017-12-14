package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"creative_info"
	"util"
)

type Conf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

var Gdb *sql.DB

func Init(cf *Conf) {
	if len(cf.Host) == 0 || len(cf.Port) == 0 || len(cf.Username) == 0 || len(cf.Database) == 0 {
		panic("no mysql host or port or username or database")
	}

	dbSrc := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&tls=skip-verify&autocommit=true",
		cf.Username, cf.Password, cf.Host, cf.Port, cf.Database)
	db, err := sql.Open("mysql", dbSrc)
	if err != nil {
		panic(err)
	}
	Gdb = db
	Gdb.SetMaxOpenConns(100)
	Gdb.SetMaxIdleConns(50)
	if err := Gdb.Ping(); err != nil {
		panic(err)
	}
}

func GetCreativeInfo(cUrl, cType string) (string, int64, error) {
	var cId string
	var cSize int64
	if err := Gdb.QueryRow("SELECT id, size FROM  creative_info WHERE url = ?", cUrl).Scan(&cId, &cSize); err != nil {
		if err != sql.ErrNoRows {
			return "", 0, err
		} else {
			cSize, err = util.GetResourceSize(cUrl, 200)
			if err != nil || cSize <= 0 {
				cSize = 0
			}
			res, err := Gdb.Exec("INSERT INTO creative_info(url, type, size) VALUES(?, ?, ?)", cUrl, cType, cSize)
			if err != nil {
				return "", 0, err
			}
			id, err := res.LastInsertId()
			if err != nil || id == 0 {
				return "", 0, err
			}
			cId = strconv.FormatInt(id, 10)
		}
	}
	return cId, cSize, nil
}

func GetCreativeInfoWithNoSize() ([]creative_info.CreativeInfo, error) {
	rows, err := Gdb.Query("SELECT id, url, fail_times FROM creative_info WHERE size=0 and fail_times<=10 ORDER BY id DESC LIMIT 50")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cInfos []creative_info.CreativeInfo
	for rows.Next() {
		var info creative_info.CreativeInfo
		if err := rows.Scan(&info.Id, &info.Url, &info.FailTimes); err != nil {
			log.Println("[db] GetCreativeInfoWithNoSize rows.Scan error: ", err)
			continue
		} else {
			cInfos = append(cInfos, info)
		}
	}

	return cInfos, nil
}

func BatchUpdateSize(cInfos []creative_info.CreativeInfo) error {
	stmt, err := Gdb.Prepare("UPDATE creative_info SET size=?, fail_times=? WHERE url=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, info := range cInfos {
		if info.Size > 0 {
			info.FailTimes = 0
		}
		if _, err := stmt.Exec(info.Size, info.FailTimes, info.Url); err != nil {
			log.Println("[db] BatchUpdateSize error: ", err, ", url: ", info.Url, ", size: ", info.Size)
			continue
		}
	}

	return nil
}
