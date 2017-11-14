package db

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type Conf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

var db *sql.DB

func Init(cf *Conf) {
	if len(cf.Host) <= 0 || len(cf.Port) <= 0 || len(cf.Username) <= 0 || len(cf.Database) <= 0 {
		panic("no mysql host or port or username or database")
	}

	_, err := strconv.Atoi(cf.Port)
	if err != nil {
		panic("mysql port not number: " + cf.Port)
	}

	dbSrc := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		cf.Username, cf.Password, cf.Host, cf.Port, cf.Database)
	db, err = sql.Open("mysql", dbSrc)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	if err := db.Ping(); err != nil {
		panic(err)
	}
}

func GetCreativeId(cUrl string) (int64, error) {
	var id string
	stmt, err := db.Prepare("SELECT id FROM  creative_info WHERE url = ?")
	if err != nil {
		return 0, err
	}
	if stmt != nil {
		defer stmt.Close()
	}
	if err = stmt.QueryRow(cUrl).Scan(&id); err != nil {
		if err != sql.ErrNoRows {
			return 0, err
		} else {
			if _, err := db.Exec("INSERT INTO creative_info(url) VALUES(?)", cUrl); err != nil {
				return 0, err
			}
			if err := stmt.QueryRow(cUrl).Scan(&id); err != nil {
				return 0, err
			}
		}
	}
	cId, _ := strconv.ParseInt(id, 10, 64)
	return cId, nil
}
