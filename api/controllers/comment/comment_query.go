package comment

import (
	"config"
	"models"
)

func getNumberOfCommentById(video string) (count int, err error) {
	err = (*config.API.Db).QueryRow("SELECT COUNT(*) FROM comment WHERE video_id = ?", video).Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func getComments(page int, ppage int) ([]models.Comment, error) {
	count := 0
	rows, err := (*config.API.Db).Query("SELECT * FROM comment ORDER BY body LIMIT ?,?;", (ppage)*(page-1), ppage)

	res := make([]models.Comment, 0)
	if err != nil {
		return res, err
	}
	defer rows.Close()
	for rows.Next() {
		var com models.Comment
		err = rows.Scan(&com.Id, &com.Body, &com.User_id, &com.Video_id)
		if err != nil {
			return res, err
		}
		res = append(res, com)
		count++
	}
	if err != nil {
		return res, err
	}
	return res, nil
}
