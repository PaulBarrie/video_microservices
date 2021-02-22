package config

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go/v7"
)

//API defines utils cli & router
var API *App = &App{}

//App Struct defining API var
type App struct {
	Router *mux.Router
	Minio  *minio.Client
}

func init() {
	API.ConnectMinio()
}

//Run the API
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
