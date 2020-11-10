package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"

	"config"
	"controllers/utils"

	"github.com/gorilla/mux"
)

//CreateVideo manages the root to create a new video
func CreateVideo(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	check5 := utils.Check500(w)
	token := r.FormValue("token")
	id := mux.Vars(r)["id"]

	// Check information sent are compliants
	check := utils.Check400M(w)
	if check(name, "name") {
		return
	}

	if !utils.IsDoubleAuthentified(token, id) {
		http.Error(w, "[401]- Unauthorized: you must be authentified !", http.StatusUnauthorized)
		return
	}
	if _, _, err := r.FormFile("file"); err != nil {
		http.Error(w, "[401]- Unauthorized: you must provide a file !", http.StatusUnauthorized)
		return
	}
	/* Insert base video */
	stmt, err := (*config.API.Db).Prepare("INSERT INTO video ( name, duration, user_id, source, created_at, view, enabled) VALUES (?, ?, ?, ?, NOW(), ?, ?)")
	if check5(err) {
		return
	}
	res, err := stmt.Exec(name, 0, id, "", 0, 1)
	if check5(err) {
		return
	}
	uid, err := res.LastInsertId()
	if check5(err) {
		return
	}
	uidStr := strconv.FormatInt(uid, 10)
	video := getVideoByID(uidStr)
	videoP, err := upload(w, r, video)

	if check5(err) {
		return
	}

	fields := map[string]string{
		"duration": fmt.Sprintf("%f", videoP.Duration),
		"source":   videoP.Source,
	}
	err = queryOnUpdateVideos(fields, uidStr)
	if check5(err) {
		return
	}
	/* Send response */
	video = getVideoByID(uidStr)
	js, _ := json.Marshal(utils.RespVideo{"ok", video})

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

//GetVideoList allows to get a list of videos
func GetVideoList(w http.ResponseWriter, r *http.Request) {
	//name := r.Header.Get("name")
	page := r.Header.Get("page")
	perPage := r.Header.Get("perPage")

	// Check information sent are compliants
	if page == "" {
		page = "1"
	}
	if perPage == "" {
		perPage = "10"
	}
	pp, _ := strconv.Atoi(perPage)
	p, _ := strconv.Atoi(page)
	video, err := queryVideos(p, pp)

	check5 := utils.Check500(w)
	if check5(err) {
		return
	}
	/* Send response */
	cmpt, err := getNumberOfVideos()
	if check5(err) {
		return
	}
	js, _ := json.Marshal(utils.RespVideoPaginated{"ok", video, utils.Pager{p, int(math.Ceil(float64(cmpt) / float64(pp)))}})

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

//GetVideoListByUser allows to get the video list for a given user
func GetVideoListByUser(w http.ResponseWriter, r *http.Request) {
	page := r.Header.Get("page")
	perPage := r.Header.Get("perPage")
	id := mux.Vars(r)["id"]

	if page == "" {
		page = "1"
	}
	if perPage == "" {
		perPage = "10"
	}
	pp, _ := strconv.Atoi(perPage)
	p, _ := strconv.Atoi(page)

	/* Make query */
	check5 := utils.Check500(w)

	video, err := getUserVideos(id, p, pp)
	/* Send response */
	cmpt, err := getNumberOfVideosByUser(id)
	if check5(err) {
		return
	}
	js, _ := json.Marshal(utils.RespVideoPaginated{"ok", video, utils.Pager{p, int(math.Ceil(float64(cmpt) / float64(pp)))}})

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

// EncodeVideoByID allows to retrieve a given video
func EncodeVideoByID(w http.ResponseWriter, r *http.Request) {
	// grab the generated receipt.pdf file and stream it to browser
	id := mux.Vars(r)["id"]
	vid := getVideoByID(id)
	streamPDFbytes, err := ioutil.ReadFile(vid.Source)
	log.Println(r)
	check5 := utils.Check500(w)
	if check5(err) {
		return
	}

	b := bytes.NewBuffer(streamPDFbytes)

	// stream straight to client(browser)
	w.Header().Set("Content-type", "video/mp4")

	if _, err := b.WriteTo(w); err != nil { // <----- here!
		fmt.Fprintf(w, "%s", err)
	}

	w.Write([]byte("Video Completed"))
}

//UpdateVideo allows to update video fields
func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	uid := mux.Vars(r)["id"]
	fields := map[string]string{
		"name":    r.FormValue("name"),
		"user_id": r.FormValue("user"),
	}
	if !utils.IsDoubleAuthentified(token, r.FormValue("user")) {
		http.Error(w, "[401]- Unauthorized: you must be authentified !", http.StatusUnauthorized)
		return
	}
	err := queryOnUpdateVideos(fields, uid)
	check5 := utils.Check500(w)
	if check5(err) {
		return
	}
	video := getVideoByID(uid)
	/* Send response */
	js, _ := json.Marshal(utils.RespVideo{"ok", video})

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

//DeleteVideo allows to delete a video
func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	id := mux.Vars(r)["id"]
	userID := ""
	err := (*config.API.Db).QueryRow("SELECT user_id FROM token WHERE code = ?", token).Scan(&userID)
	check5 := utils.Check500(w)
	if check5(err) {
		return
	}

	if !utils.IsDoubleAuthentified(token, userID) {
		http.Error(w, "[401]- Unauthorized: you must be authentified !", http.StatusUnauthorized)
		return
	}
	_, err = (*config.API.Db).Exec("DELETE FROM video WHERE id=?", id)
	if check5(err) {
		return
	}
	//Delete video from sys file
	w.WriteHeader(http.StatusCreated)
}
