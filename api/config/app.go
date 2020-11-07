package config

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go/v7"
)

//Api defines utils cli & router
var Api *App = &App{}

//App Struct defining Api var
type App struct {
	Router  *mux.Router
	Db      *sql.DB
	Minio   *minio.Client
	Postfix *PostfixCli
}

func init() {
	Api.InitDB()
	Api.ConnectMinio()
	Api.InitPostfix()
}

//Run the API
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
