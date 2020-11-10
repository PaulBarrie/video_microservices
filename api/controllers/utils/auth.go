package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"config"
	"models"
)

//Authentify defines the route for yser auth
func Authentify(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	// Check information sent are compliants
	if login == "" {
		http.Error(w, "[400]- Bad request: login is missing", http.StatusBadRequest)
		return
	}
	if password == "" {
		http.Error(w, "[400]- Bad request: password is missing", http.StatusBadRequest)
		return
	}
	userID := checkLoginDetails(login, password)
	if userID == -1 {
		http.Error(w, "[403]- Forbidden: wrong login details", http.StatusForbidden)
		return
	} else if id := userHasToken(userID); id != -1 {
		err := removeToken(userID)
		if err != nil {
			http.Error(w, "[500]- Error: internal server error", http.StatusInternalServerError)
			return
		}
	}

	token := createToken(userID)
	registerToken(token)
	tokenReg, err := getTokenByUserID(userID)

	if err != nil {
		http.Error(w, "[500]- Error: internal server error", http.StatusInternalServerError)
		return
	}
	js, _ := json.Marshal(RespToken{"ok", tokenReg})
	w.WriteHeader(http.StatusCreated)
	w.Write(js)

}

//IsDoubleAuthentified allows to identify user with token and id
func IsDoubleAuthentified(token string, idIn string) bool {
	var id string
	sqlStatement := `SELECT user_id FROM token WHERE code = ?;`
	row := (*config.API.Db).QueryRow(sqlStatement, token)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	} else if idIn == id {
		return true
	}
	return false
}

//IsAuthentified allows "to identify user with token onlt
func IsAuthentified(token string) bool {
	var id string
	sqlStatement := `SELECT user_id FROM token WHERE code = ?;`
	row := (*config.API.Db).QueryRow(sqlStatement, token)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	return true
}

func checkLoginDetails(login string, pwd string) int {
	stmt1 := `SELECT id FROM user WHERE email = ? AND password = ?;`
	stmt2 := `SELECT id FROM user WHERE username = ? AND password = ?;`
	pwd = HashPwd(pwd)
	check1 := checkAndGetUserID(stmt1, login, pwd)
	if check1 > -1 {
		return check1
	}
	return checkAndGetUserID(stmt2, login, pwd)
}

func checkAndGetUserID(stmt string, login string, pwd string) int {
	var id int
	row := (*config.API.Db).QueryRow(stmt, login, pwd)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1
		}
	}

	return id
}

func createToken(id int) models.Token {
	code := make([]byte, 8)
	rand.Read(code)

	return models.Token{
		Id:         -1,
		Code:       fmt.Sprintf("%x", code),
		Expired_at: (time.Now()).AddDate(0, 0, 2), // Expire in 2 days
		User_id:    id}
}

/* Queries on token table */

func userHasToken(userID int) int {
	var id int
	sqlStatement := `SELECT id FROM token WHERE user_id = ?;`

	row := (*config.API.Db).QueryRow(sqlStatement, userID)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1
		}
	}
	return id
}

func removeToken(userID int) error {
	del, err := (*config.API.Db).Prepare("DELETE FROM token WHERE user_id=?;")
	if err != nil {
		return err
	}
	del.Exec(userID)
	return nil
}

func getTokenByUserID(userID int) (models.Token, error) {
	stmt := `SELECT * FROM token WHERE user_id = ?;`
	token := models.Token{}

	row := (*config.API.Db).QueryRow(stmt, userID)
	err := row.Scan(&token.Id, &token.Code, &token.Expired_at, &token.User_id)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Token{}, err
		}
	}
	return token, nil
}

func registerToken(token models.Token) {
	query := fmt.Sprintf(
		"INSERT INTO token (code, expired_at, user_id) VALUES ('%s', '%s', '%s');",
		token.Code, (token.Expired_at).Format("2006-01-02 15:04:05"), strconv.Itoa(token.User_id))
	insert, err := (*config.API.Db).Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

//HashPwd returns a sh256 hash of pwd
func HashPwd(pwd string) string {
	hashPwd := sha256.Sum256([]byte(pwd))

	return hex.EncodeToString(hashPwd[:])
}
