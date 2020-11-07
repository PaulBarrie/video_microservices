package config

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go/v7"
)

//Api defines utils cli & router
var Api *App = &App{}

//App Struct defining Api var
type App struct {
	Router *mux.Router
	Minio  *minio.Client
}

func init() {
	Api.ConnectMinio()
}

//Run the API
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
