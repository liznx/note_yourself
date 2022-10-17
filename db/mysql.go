package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var MysqlDB *sql.DB

func init() {
	var err error
	MysqlDB, err = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/laizhongxiao?charset=utf8mb4")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = MysqlDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}
