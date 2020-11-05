package router

import (
	"controllers/comment"
	"controllers/user"
	"controllers/utils"
	"controllers/video"

	"github.com/gorilla/mux"
)

//InitializeRouter init a router instance
func InitializeRouter() *mux.Router {
	router := mux.NewRouter() //.StrictSlash(true)

	/*
		ROUTES FOR SERVICES
	*/

	/* ROUTES FOR USERS */
	router.HandleFunc("/user", user.RegisterUser).Methods("POST")
	router.HandleFunc("/auth", utils.Authentify).Methods("POST")
	router.HandleFunc("/user/{id}", user.DeleteUser).Methods("DELETE")
	router.HandleFunc("/user/{id}", user.UpdateUser).Methods("PUT")
	router.HandleFunc("/users", user.GetUsers).Methods("GET")
	router.HandleFunc("/user/{id}", user.GetUserById).Methods("GET")
	/* ROUTES FOR VIDEO */
	router.HandleFunc("/user/{id}/video", video.CreateVideo).Methods("POST")
	router.HandleFunc("/videos", video.GetVideoList).Methods("GET")
	router.HandleFunc("/user/{id}/videos", video.GetVideoListByUser).Methods("GET")
	router.HandleFunc("/video/{id}", video.EncodeVideoByID).Methods("PATCH")
	router.HandleFunc("/video/{id}", video.UpdateVideo).Methods("PUT")
	router.HandleFunc("/video/{id}", video.DeleteVideo).Methods("DELETE")
	/* ROUTES FOR COMMENTARY */
	router.HandleFunc("/video/{id}/comment", comment.CreateComment).Methods("POST")
	router.HandleFunc("/video/{id}/comments", comment.GetCommentsList).Methods("GET")
	// Documentation
	//router.PathPrefix("/docu").Handler(httpSwagger.WrapHandler)

	return router
}
