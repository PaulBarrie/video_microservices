package config

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go/v7"
)

//Api defines utils cli & router
var API *App = &App{}

//App Struct defining API var
type App struct {
	Router *mux.Router
	Db     *sql.DB
	Minio  *minio.Client
	Smtp   *SmtpCli
}

func init() {
	API.InitDB()
	API.ConnectMinio()
	API.InitSMTP()
}

//Run the API
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
