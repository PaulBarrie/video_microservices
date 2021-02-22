package user

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"config"
	"controllers/utils"

	"github.com/gorilla/mux"
)

// YoutubeAPI godoc
// @Summary Create user
// @Description Get details of all orders
// @Tags orders
// @Accept  json
// @Produce  json
// @Success 200
// @Router /user [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	uname := r.FormValue("username")
	pseudo := r.FormValue("pseudo")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if uname == "" {
		http.Error(w, "[400]- Bad request: username is missing", http.StatusBadRequest)
		return
	}
	if email == "" {
		http.Error(w, "[400]- Bad request: email is missing", http.StatusBadRequest)
		return
	}
	if password == "" {
		http.Error(w, "[400]- Bad request: password is missing", http.StatusBadRequest)
		return
	}
	if !isEmailAvailable(email) {
		http.Error(w, "[403]- Forbidden: email is already used", http.StatusBadRequest)
		return
	}
	if !isUnameAvailable(uname) {
		http.Error(w, "[403]- Forbidden: username is already used", http.StatusBadRequest)
		return
	}
	/* Make query */
	query := fmt.Sprintf(
		"INSERT INTO user ( username, email, pseudo, password, created_at) VALUES ('%s', '%s', '%s', '%s', NOW());",
		uname, email, pseudo, utils.HashPwd(password))
	insert, err := (*config.API.Db).Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
	user := getUserByEmail(email)
	/* Send response */

	js, _ := json.Marshal(utils.RespUser{"ok", user})
	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

//DeleteUser allows to delete account
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	id := mux.Vars(r)["id"]

	if !utils.IsDoubleAuthentified(token, id) {
		http.Error(w, "[401]- Unauthorized: you must be authentified !", http.StatusUnauthorized)
		return
	}
	_, err := (*config.API.Db).Exec("DELETE FROM user WHERE id=?", id)
	if err != nil {
		http.Error(w, "[500]- Error: internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	id := mux.Vars(r)["id"]
	fields := map[string]string{
		"username": r.FormValue("username"),
		"pseudo":   r.FormValue("pseudo"),
		"email":    r.FormValue("email"),
		"password": r.FormValue("password"),
	}
	if !utils.IsDoubleAuthentified(token, id) {
		http.Error(w, "[401]- Unauthorized: you must be authentified !", http.StatusUnauthorized)
		return
	}
	err := queryOnUpdate(fields, id)
	if err != nil {
		http.Error(w, "[500]- Error: internal server error", http.StatusInternalServerError)
		return
	}
	user := ReqUserByID(id)
	/* Send response */
	js, _ := json.Marshal(utils.RespUser{"ok", user})

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	page := r.Header.Get("page")
	p_page := r.Header.Get("perPage")

	cmpt, err := getNumberOfUsers()
	check5 := utils.Check500(w)
	if check5(err) {
		return
	}
	if page == "" {
		page = "1"
	}
	if p_page == "" {
		p_page = "10"
	}
	pp, _ := strconv.Atoi(p_page)
	p, _ := strconv.Atoi(page)
	user, err := queryUsers(p, pp)
	/* Send response */

	js, _ := json.Marshal(utils.RespUserPaginated{"ok", user, utils.Pager{p, int(math.Ceil(float64(cmpt) / float64(pp)))}})

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

//GetUserByID handles the route to retriev an user with a given ID.
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	id := mux.Vars(r)["id"]

	if !utils.IsAuthentified(token) {
		http.Error(w, "[401]- Unauthorized: you must be authentified !", http.StatusUnauthorized)
		return
	}
	/* Make query */
	user := ReqUserByID(id)
	/* Send response */
	js, _ := json.Marshal(utils.RespUser{"ok", user})

	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
