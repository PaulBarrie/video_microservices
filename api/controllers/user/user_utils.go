package user

import (
	"database/sql"

	"config"
	"models"
)

/* Utils function */

func scanUserRow(stmt string, val string) models.User {
	usr := models.User{}

	row := (*config.Api.Db).QueryRow(stmt, val)
	err := row.Scan(&usr.Id, &usr.Username, &usr.Email, &usr.Pseudo, &usr.Password, &usr.Created_at)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}
		} else {
			panic(err)
		}
	}

	return usr
}

func isResultEmpty(stmt string, val string) bool {
	var id int
	row := (*config.Api.Db).QueryRow(stmt, val)
	err := row.Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return true
		} else {
			panic(err)
		}
	}

	return false
}
