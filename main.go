package main

import (
	"noteyouself/db"
)

func main() {
	defer db.MysqlDB.Close()
	router := initRouter()
	router.Run(":8001")
}
