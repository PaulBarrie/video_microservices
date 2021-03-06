package video

import (
	"config"
	"models"
)

/* Query utility functions */
func getNumberOfVideos() (count int, err error) {
	err = (*config.API.Db).QueryRow("SELECT COUNT(*) FROM video").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func getNumberOfVideosByUser(usr string) (count int, err error) {
	err = (*config.API.Db).QueryRow("SELECT COUNT(*) FROM video WHERE user_id = ?", usr).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func queryVideos(page int, ppage int) ([]models.Video, error) {
	count := 0
	rows, err := (*config.API.Db).Query("SELECT * FROM video ORDER BY name LIMIT ?,?;", (ppage)*(page-1), ppage)

	res := make([]models.Video, 0)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var vid models.Video
		err = rows.Scan(&vid.Id, &vid.Name, &vid.Duration, &vid.User_id, &vid.Source, &vid.Created_at, &vid.View, &vid.Enabled)
		if err != nil {
			return res, err
		}
		res = append(res, vid)
		count++
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

func queryOnUpdateVideos(fields map[string]string, id string) error {
	qSet := ""
	for key, val := range fields {
		if val != "" {
			qSet += key + " = '" + val + "', "
		}
	}
	qSet = qSet[:len(qSet)-2]
	stmt := "UPDATE video SET " + qSet + " WHERE id = ?;"
	_, err := (*config.API.Db).Exec(stmt, id)

	return err
}

func getUserVideos(uid string, page int, ppage int) ([]models.Video, error) {
	rows, err := (*config.API.Db).Query("SELECT * FROM video WHERE user_id = ? AND enabled = 1 ORDER BY created_at DESC LIMIT ?,?;", uid, (ppage)*(page-1), ppage)
	res := make([]models.Video, 0)
	count := 0
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var vid models.Video
		err = rows.Scan(&vid.Id, &vid.Name, &vid.Duration, &vid.User_id, &vid.Source, &vid.Created_at, &vid.View, &vid.Enabled)
		if err != nil {
			return res, err
		}
		res = append(res, vid)
		count++
	}
	if err != nil {
		return res, err
	}
	return res, nil
}

func getVideoByID(id string) models.Video {
	sqlStatement := `SELECT * FROM video WHERE id = ?;`
	return scanVideoRow(sqlStatement, id)
}
