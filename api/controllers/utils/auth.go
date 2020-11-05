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
	user_id := checkLoginDetails(login, password)
	if user_id == -1 {
		http.Error(w, "[403]- Forbidden: wrong login details", http.StatusForbidden)
		return
	} else if id := userHasToken(user_id); id != -1 {
		err := removeToken(user_id)
		if err != nil {
			http.Error(w, "[500]- Error: internal server error", http.StatusInternalServerError)
			return
		}
	}

	token := createToken(user_id)
	registerToken(token)
	token_reg, err := getTokenByUserId(user_id)

	if err != nil {
		http.Error(w, "[500]- Error: internal server error", http.StatusInternalServerError)
		return
	}
	js, _ := json.Marshal(RespToken{"ok", token_reg})
	w.WriteHeader(http.StatusCreated)
	w.Write(js)

}

/* Auth utility */
func IsDoubleAuthentified(token string, id_in string) bool {
	var id string
	sqlStatement := `SELECT user_id FROM token WHERE code = ?;`
	row := (*config.Api.Db).QueryRow(sqlStatement, token)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err)
		}
	} else if id_in == id {
		return true
	} else {
		return false
	}
}
func IsAuthentified(token string) bool {
	var id string
	sqlStatement := `SELECT user_id FROM token WHERE code = ?;`
	row := (*config.Api.Db).QueryRow(sqlStatement, token)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			panic(err)
		}
	}

	return true

}

func checkLoginDetails(login string, pwd string) int {
	stmt1 := `SELECT id FROM user WHERE email = ? AND password = ?;`
	stmt2 := `SELECT id FROM user WHERE username = ? AND password = ?;`
	pwd = HashPwd(pwd)
	check1 := checkAndGetUserId(stmt1, login, pwd)
	if check1 > -1 {
		return check1
	} else {
		return checkAndGetUserId(stmt2, login, pwd)
	}
}

func checkAndGetUserId(stmt string, login string, pwd string) int {
	var id int
	row := (*config.Api.Db).QueryRow(stmt, login, pwd)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1
		} else {
			panic(err)
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

func userHasToken(user_id int) int {
	var id int
	sqlStatement := `SELECT id FROM token WHERE user_id = ?;`

	row := (*config.Api.Db).QueryRow(sqlStatement, user_id)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return -1
		} else {
			panic(err)
		}
	}
	return id
}

func removeToken(user_id int) error {
	del, err := (*config.Api.Db).Prepare("DELETE FROM token WHERE user_id=?;")
	if err != nil {
		return err
	}
	del.Exec(user_id)
	return nil
}

func getTokenByUserId(usr_id int) (models.Token, error) {
	stmt := `SELECT * FROM token WHERE user_id = ?;`
	token := models.Token{}

	row := (*config.Api.Db).QueryRow(stmt, usr_id)
	err := row.Scan(&token.Id, &token.Code, &token.Expired_at, &token.User_id)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Token{}, err
		} else {
			return models.Token{}, err
		}
	}

	return token, nil
}

func registerToken(token models.Token) {
	query := fmt.Sprintf(
		"INSERT INTO token (code, expired_at, user_id) VALUES ('%s', '%s', '%s');",
		token.Code, (token.Expired_at).Format("2006-01-02 15:04:05"), strconv.Itoa(token.User_id))
	insert, err := (*config.Api.Db).Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
}

func HashPwd(pwd string) string {
	hash_pwd := sha256.Sum256([]byte(pwd))

	return hex.EncodeToString(hash_pwd[:])
}
