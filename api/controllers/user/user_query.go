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

func getUserById(id string) models.User {
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
	q_set := ""
	empty_fields := 0
	for key, val := range fields {
		if val != "" {
			q_set += key + " = '" + val + "', "
		} else {
			empty_fields++
		}
		if empty_fields == 4 {
			return nil
		}
	}
	q_set = q_set[:len(q_set)-2]
	stmt := "UPDATE user SET " + q_set + " WHERE id = ?;"
	fmt.Println(stmt)
	_, err := (*config.Api.Db).Exec(stmt, id)

	return err
}

/* Query several users */

func queryUsers(page int, ppage int) ([]models.User, error) {
	count := 0
	rows, err := (*config.Api.Db).Query("SELECT * FROM user ORDER BY username LIMIT ?,?;", (ppage)*(page-1), ppage)

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
	err = (*config.Api.Db).QueryRow("SELECT COUNT(*) FROM user").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}
