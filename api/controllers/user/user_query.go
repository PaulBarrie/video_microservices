package user

import (
	"fmt"

	"config"
	"models"
)

/* Test unicity of informations */

func isEmailAvailable(email string) bool {
	sqlStatement := `SELECT id FROM user WHERE email = ?;`
	return isResultEmpty(sqlStatement, email)
}

func isUnameAvailable(uname string) bool {
	sqlStatement := `SELECT id FROM user WHERE username = ?;`
	return isResultEmpty(sqlStatement, uname)
}

/* Query single row */

//ReqUserByID retrieves a user in DB with a given ID.
func ReqUserByID(id string) models.User {
	sqlStatement := `SELECT * FROM user WHERE id = ?;`

	return scanUserRow(sqlStatement, id)
}

func getUserByEmail(email string) models.User {
	sqlStatement := `SELECT * FROM user WHERE email = ?;`

	return scanUserRow(sqlStatement, email)
}

func getUserByUname(uname string) models.User {
	sqlStatement := `SELECT * FROM user WHERE username = ?;`

	return scanUserRow(sqlStatement, uname)
}

func getUserByPseudo(uname string) models.User {
	sqlStatement := `SELECT * FROM user WHERE pseudo = ?;`

	return scanUserRow(sqlStatement, uname)
}

/* Query for update */

func queryOnUpdate(fields map[string]string, id string) error {
	qSet := ""
	emptyFields := 0
	for key, val := range fields {
		if val != "" {
			qSet += key + " = '" + val + "', "
		} else {
			emptyFields++
		}
		if emptyFields == 4 {
			return nil
		}
	}
	qSet = qSet[:len(qSet)-2]
	stmt := "UPDATE user SET " + qSet + " WHERE id = ?;"
	fmt.Println(stmt)
	_, err := (*config.API.Db).Exec(stmt, id)

	return err
}

/* Query several users */

func queryUsers(page int, ppage int) ([]models.User, error) {
	count := 0
	rows, err := (*config.API.Db).Query("SELECT * FROM user ORDER BY username LIMIT ?,?;", (ppage)*(page-1), ppage)

	res := make([]models.User, 0)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var usr models.User
		err = rows.Scan(&usr.Id, &usr.Username, &usr.Email, &usr.Pseudo, &usr.Password, &usr.Created_at)
		if err != nil {
			return res, err
		}
		res = append(res, usr)
		count++
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

func getNumberOfUsers() (count int, err error) {
	err = (*config.API.Db).QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}
