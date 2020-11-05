package video

import (
	"encoding/json"
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
	//video.Duration = header.Size
	// Copy the file data to my buffer
	io.Copy(&buf, file)
	defer buf.Reset()
	tmpLoc := fmt.Sprintf("tmp/%s", video.Source)
	err = saveFile(&buf, tmpLoc)
	if err != nil {
		return video, err
	}
	defer os.Remove(tmpLoc)
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
	err = saveInMinio(tmpFile, path, video.Name, details.Quality)
	if err != nil {
		return video, err
	}
	err = go sendEncodeRequest(path, details.Quality, video.Name)
	if err != nil {
		return video, err
	}
	return video, nil
}

func saveFile(buf *bytes.Buffer, path string) error {
	//file, err := os.OpenFile(path, os.O_CREATE, 0666) //os.O_WRONLY|os.O_TRUNC|
	p_path := strings.Split(path, "/")
	p_pathjoin := strings.Join(p_path[:len(p_path)-1], "/")
	if _, err := os.Stat(p_pathjoin); os.IsNotExist(err) {
		err = os.Mkdir(p_pathjoin, 0755)
		if err != nil {
			return err
		}
	}

	err := ioutil.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func getVideoDetails(path string) (utils.VideoDetails, error) {
	re := regexp.MustCompile("[0-9]+")
	// Get quality
	cmd_q := fmt.Sprintf("ffprobe -v error -show_format -show_streams %s | grep pix_fmt", path)
	out_q, err := exec.Command("sh", "-c", cmd_q).Output()
	if err != nil {
		return utils.VideoDetails{0, 0}, err
	}
	quality := re.FindAllString(string(out_q), -1)[0]
	q_int, _ := strconv.Atoi(quality)
	//Get duration
	cmd_d := fmt.Sprintf("ffprobe -v error -show_format -show_streams %s | grep duration | tail -1", path)
	out_d, err := exec.Command("sh", "-c", cmd_d).Output()
	if err != nil {
		return utils.VideoDetails{0, 0}, err
	}
	duration := strings.Split(string(out_d), "=")[1]
	duration = duration[:len(duration)-1]
	d_int, err := strconv.ParseFloat(duration, 64)
	if err != nil {
		return utils.VideoDetails{0, 0}, err
	}
	return utils.VideoDetails{d_int, int64(q_int)}, nil
}

func saveInMinio(file *os.File, bucket_name string, fileName string, fileResolution int64) error {
	cli := config.Api.Minio
	log.Println(bucket_name)
	err := cli.MakeBucket(context.Background(), bucket_name, minio.MakeBucketOptions{ObjectLocking: false})
	if err != nil {
		log.Println("Errro in making bucket")
		log.Println(err)
		return err
	}
	log.Println("Bucket created")
	fileStat, err := file.Stat()
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = cli.PutObject(context.Background(), bucket_name, fmt.Sprintf("%d_%s", fileResolution, fileName), file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func makeBucketName(name string, idUsr int, idVid int, createdAt string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	createdAt = createdAt[:len(createdAt)-3]
	return fmt.Sprintf("usr%d%s%d%s", idUsr, name, idVid, reg.ReplaceAllString(createdAt, "")), nil

}

func sendEncodeRequest(bucketName string, format int64, filename string) error {
	log.Println("Sending encode request")

	formData := url.Values{
		"bucket_name": {bucketName},
		"format": {format},
		"filename": {filename}
	}
	if err != nil {
		log.Println(err)
		return err
	}
	defer reqBody.Close()

	resp, err := http.PostForm("http://video_encoder/encode", formData)
	if err != nil {
		log.Println(err)
		return err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	log.Println(result["form"])
	return nil
}

/* Utils for query */
func scanVideoRow(stmt string, val string) models.Video {
	vid := models.Video{}

	row := (*config.Api.Db).QueryRow(stmt, val)
	err := row.Scan(&vid.Id, &vid.Name, &vid.Duration, &vid.User_id, &vid.Source, &vid.Created_at, &vid.View, &vid.Enabled)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Video{}
		} else {
			panic(err)
		}
	}

	return vid
}

func idOrName(user string, w http.ResponseWriter) (string, error) {
	if _, err := strconv.Atoi(user); err == nil {
		var id string
		st := fmt.Sprintf(
			"SELECT id FROM user WHERE username = '%s' OR email = '%s'",
			user, user)
		row, err := (*config.Api.Db).Query(st)
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
	} else {
		return user, nil
	}
}
