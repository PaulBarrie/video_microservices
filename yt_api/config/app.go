package config

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go/v7"
)

var Api *App = &App{}

type App struct {
	Router *mux.Router
	Db     *sql.DB
	Minio  *minio.Client
}

func init() {
	Api.InitDB()
	Api.ConnectMinio()
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
