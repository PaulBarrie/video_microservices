package video

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"log"
	"strconv"
	"context"
	"strings"
	minio "github.com/minio/minio-go/v7"
	"config"
	"models"
	"controllers/utils"
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
	tmp_loc := fmt.Sprintf("tmp/%s", video.Source)
	err = saveFile(&buf, tmp_loc)
	if err != nil {
		return video, err
	}
	defer os.Remove(tmp_loc)
	details, err := getVideoDetails(tmp_loc)
	if err != nil {
		return video, err
	}
	video.Duration = details.Duration
	log.Println(details)
	tmp_file, err :=  os.OpenFile(tmp_loc, os.O_CREATE, 0755)
	if err != nil {
		return video, err
	}
	defer os.Remove(tmp_loc)
	err = saveInMinio(tmp_file, path, video.Name, details.Quality)
	if err != nil {
		return video, err
	}
	err = sendEncodeRequest(path,details.Quality, video.Name)
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
	out_q, err := exec.Command("sh","-c",cmd_q).Output()
	if err != nil {
        return utils.VideoDetails{0,0}, err
	}
	quality := re.FindAllString(string(out_q), -1)[0]
	q_int,_ := strconv.Atoi(quality)
	//Get duration
	cmd_d := fmt.Sprintf("ffprobe -v error -show_format -show_streams %s | grep duration | tail -1", path)
	out_d, err := exec.Command("sh","-c",cmd_d).Output()
	if err != nil {
        return utils.VideoDetails{0,0}, err
	}
	duration := strings.Split(string(out_d), "=")[1]
	duration = duration[: len(duration)-1]
	d_int, err := strconv.ParseFloat(duration, 64)
	if err != nil {
        return utils.VideoDetails{0,0}, err
	}
	return utils.VideoDetails{d_int,int64(q_int)}, nil
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
	info, err := cli.PutObject(context.Background(), bucket_name, fmt.Sprintf("%d_%s", fileResolution, fileName), file, fileStat.Size(), minio.PutObjectOptions{ContentType:"application/octet-stream"})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(info)
	return nil
}


func makeBucketName(name string, id_usr int, id_vid int, created_at string) (string, error) {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	created_at = created_at[:len(created_at)-3]
	return fmt.Sprintf("usr%d%s%d%s", id_usr, name, id_vid, reg.ReplaceAllString(created_at, "")), nil

}

func sendEncodeRequest(bucket_name string, format int64, filename string) error {
	log.Println("Sending encode request")
	url := "http://video_encoder/encode"
	var body = []byte(fmt.Sprintf(`{"bucket_name": %s, "format": "%d", "filename": %s}`, bucket_name, format, filename))
	cli := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	log.Println("response Status:", resp.Status)
	resp_body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(resp_body))
	// defer resp.Body.Close()
	// ch := make(chan error)
	// go func(url string) error{
	// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	// 	resp, err := cli.Do(req)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Println("response Status:", resp.Status)
	// 	resp_body, _ := ioutil.ReadAll(resp.Body)
	// 	log.Println("response Body:", string(resp_body))
	// 	defer resp.Body.Close()
	// 	return nil
	// }(url)
	// x := <- ch
	// return x
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
