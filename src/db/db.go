package db

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

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
	if err := Gdb.QueryRow("SELECT id, size FROM  creative_info WHERE url = ?", cUrl).Scan(&cId, cSize); err != nil {
		if err != sql.ErrNoRows {
			return "", 0, err
		} else {
			cSize, err = util.GetResourceSize(cUrl)
			if err != nil || cSize <= 0 {
				return "", 0, err
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
