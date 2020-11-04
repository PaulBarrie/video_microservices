package comment

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"config"
	"models"
	"controllers/utils"
	"github.com/gorilla/mux"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	body := r.FormValue("body")
	token := r.FormValue("token")
	uid := r.FormValue("user")
	vid := mux.Vars(r)["id"]
	var cid int
	// Check information sent are compliants
	check := utils.Check400M(w)
	if check(body, "body") {
		return
	}

	if !utils.IsDoubleAuthentified(token, uid) {
		http.Error(w, "[401]- Unauthorized: you must be authentified !", http.StatusUnauthorized)
		return
	}

	/* Make query */
	row := (*config.Api.Db).QueryRow("INSERT INTO comment ( body, user_id, video_id) VALUES (?, ?, ?) SELECT LAST_INSERT_ID()", body, uid, vid)
	row.Scan(&cid)
	/* Send response */
	user_id, _ := strconv.Atoi(uid)
	video_id, _ := strconv.Atoi(vid)
	js, _ := json.Marshal(utils.RespComment{"ok", models.Comment{1, body, user_id, video_id}})

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

func GetCommentsList(w http.ResponseWriter, r *http.Request) {
	page := r.Header.Get("page")
	perPage := r.Header.Get("perPage")
	token := r.FormValue("token")
	user_id := r.FormValue("user")
	video_id := mux.Vars(r)["id"]
	if page == "" {
		page = "1"
	}
	if perPage == "" {
		perPage = "10"
	}
	p, _ := strconv.Atoi(page)
	pp, _ := strconv.Atoi(perPage)

	comment, err := getComments(p, pp)

	if err != nil {
		http.Error(w, "[500]- Error: internal server error", http.StatusInternalServerError)
		return
	}
	// Check information sent are compliants

	if !utils.IsDoubleAuthentified(token, user_id) {
		http.Error(w, "[401]- Unauthorized: you must be authentified !", http.StatusUnauthorized)
		return
	}
	stmt := fmt.Sprintf("SELECT * FROM comment WHERE video_id = '%s';", video_id)
	query, err := (*config.Api.Db).Query(stmt)
	check5 := utils.Check500(w)
	if check5(err) {
		return
	}
	defer query.Close()
	/* Send response */
	cmpt, err := getNumberOfCommentById(video_id)
	if check5(err) {
		return
	}

	js, _ := json.Marshal(utils.RespCommentPaginated{"OK", comment, utils.Pager{p, int(math.Ceil(float64(cmpt) / float64(pp)))}})

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}
