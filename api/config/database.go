package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func (a *App) InitDB() {

	var user string = os.Getenv("DB_ROOT")
	var pwd string = os.Getenv("DB_ROOT_PWD")
	var db_name string = os.Getenv("DB_NAME")
	var addr string = os.Getenv("DB_ADDRESS")
	//var port string = os.Getenv("DB_PORT")

	var uri string = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, pwd, addr /*port,*/, db_name)
	// var uri string = fmt.Sprintf("%s:%s@tcp(127.0.0.1:33060)/%s?charset=utf8&parseTime=True&loc=Local", user, pwd/*, addr port,*/, db_name)

	//fmt.Printf(uri)
	db, err := sql.Open("mysql", uri)

	if err != nil {
		log.Fatalf("[-] Error while trying to connect database\n")
	}
	if err := db.Ping(); err != nil {
		panic(err.Error())
	}

	a.Db = db
	if err != nil {
		panic(err.Error())
	}
}
