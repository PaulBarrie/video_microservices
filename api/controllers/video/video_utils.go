package video

import (
	"bytes"
	"config"
	"context"
	"controllers/utils"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"models"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	minio "github.com/minio/minio-go/v7"
)

func upload(w http.ResponseWriter, r *http.Request, video models.Video) (models.Video, error) {
	r.ParseMultipartForm(32 << 20) // limit your max input length!
	var buf bytes.Buffer
	// in your case file would be fileupload
	file, _, err := r.FormFile("file")
	if err != nil {
		return video, err
	}
	defer file.Close()
	path, err := makeBucketName(video.Name, video.User_id, video.Id, (video.Created_at).String())
	if err != nil {
		return video, err
	}
	video.Source = path
	// Copy the file data to my buffer
	io.Copy(&buf, file)
	defer buf.Reset()
	tmpLoc := fmt.Sprintf("tmp/%s", video.Source)
	err = saveFile(&buf, tmpLoc)
	if err != nil {
		return video, err
	}
	details, err := getVideoDetails(tmpLoc)
	if err != nil {
		return video, err
	}
	video.Duration = details.Duration
	tmpFile, err := os.OpenFile(tmpLoc, os.O_CREATE, 0755)
	if err != nil {
		return video, err
	}
	defer os.Remove(tmpLoc)
	sFile, err := saveInMinio(tmpFile, path, video.Source, details.Quality)
	if err != nil {
		return video, err
	}
	// Remove temporary files from
	err = os.RemoveAll("/tmp")
	if err != nil {
		log.Fatal(err)
	}
	go sendEncodeRequest(path, details.Quality, sFile, r)
	if err != nil {
		return video, err
	}
	return video, nil
}

func saveFile(buf *bytes.Buffer, path string) error {
	//file, err := os.OpenFile(path, os.O_CREATE, 0666) //os.O_WRONLY|os.O_TRUNC|
	err := os.MkdirAll("/tmp", 0777)
	if err != nil {
		return err
	}
	pPath := strings.Split(path, "/")
	pPathjoin := strings.Join(pPath[:len(pPath)-1], "/")

	if _, err := os.Stat(pPathjoin); os.IsNotExist(err) {
		err = os.Mkdir(pPathjoin, 0755)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func getVideoDetails(path string) (utils.VideoDetails, error) {
	re := regexp.MustCompile("[0-9]+")
	// Get quality
	cmdQ := fmt.Sprintf("ffprobe -v error -show_format -show_streams %s | grep pix_fmt", path)
	outQ, err := exec.Command("sh", "-c", cmdQ).Output()
	if err != nil {
		return utils.VideoDetails{0, 0}, err
	}
	quality := re.FindAllString(string(outQ), -1)[0]
	qInt, _ := strconv.Atoi(quality)
	//Get duration
	cmdD := fmt.Sprintf("ffprobe -v error -show_format -show_streams %s | grep duration | tail -1", path)
	outD, err := exec.Command("sh", "-c", cmdD).Output()
	if err != nil {
		return utils.VideoDetails{0, 0}, err
	}
	duration := strings.Split(string(outD), "=")[1]
	duration = duration[:len(duration)-1]
	dInt, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		return utils.VideoDetails{0, 0}, err
	}
	return utils.VideoDetails{dInt, int64(qInt)}, nil
}

func saveInMinio(file *os.File, bucketName string, fileName string, fileResolution int64) (string, error) {
	cli := config.API.Minio
	err := cli.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{ObjectLocking: false})
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println("Bucket created")
	fileStat, err := file.Stat()
	if err != nil {
		log.Println(err)
		return "", err
	}
	savedFile := fmt.Sprintf("%d_%s", fileResolution, fileName)
	log.Printf("Filename: %s", savedFile)
	_, err = cli.PutObject(context.Background(), bucketName, savedFile, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println(err)
		return "", err
	}
	return savedFile, nil
}

func makeBucketName(name string, idUsr int, idVid int, createdAt string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	createdAt = createdAt[:len(createdAt)-3]
	return fmt.Sprintf("usr%d%s%d%s.mp4", idUsr, name, idVid, reg.ReplaceAllString(createdAt, "")), nil

}

/* Utils for query */
func scanVideoRow(stmt string, val string) models.Video {
	vid := models.Video{}

	row := (*config.API.Db).QueryRow(stmt, val)
	err := row.Scan(&vid.Id, &vid.Name, &vid.Duration, &vid.User_id, &vid.Source, &vid.Created_at, &vid.View, &vid.Enabled)
	if err != nil && err == sql.ErrNoRows {
		return models.Video{}
	}
	return vid
}

func idOrName(user string, w http.ResponseWriter) (string, error) {
	if _, err := strconv.Atoi(user); err == nil {
		var id string
		st := fmt.Sprintf(
			"SELECT id FROM user WHERE username = '%s' OR email = '%s'",
			user, user)
		row, err := (*config.API.Db).Query(st)
		if err != nil {
			return "", err
		}
		defer row.Close()
		err = row.Scan(&id)
		if err != nil {
			http.Error(w, "[500]- Error: internal server error", http.StatusInternalServerError)
			return "", err
		}
		return id, nil
	}
	return user, nil
}
