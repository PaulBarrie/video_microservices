package config

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go/v7"
)

var Api *App = &App{}

type App struct {
	Router *mux.Router
	Minio  *minio.Client
}

func init() {
	Api.ConnectMinio()
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
