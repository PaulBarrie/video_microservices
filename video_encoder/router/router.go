package router


import (
	"controller"
	"github.com/gorilla/mux"
)

func InitializeRouter() *mux.Router {
	router := mux.NewRouter() //.StrictSlash(true)

	/*
		ROUTES FOR SERVICES
	*/
	
	/* ROUTES FOR USERS */
	router.HandleFunc("/encode", controller.EncodeVideo).Methods("POST")
	return router
}
