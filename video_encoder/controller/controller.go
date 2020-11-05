package controller

import (
	"config"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	minio "github.com/minio/minio-go/v7"
)

type response struct {
	message string `json:"message"`
}

// EncodeVideo allows to encode video in Minio
func EncodeVideo(w http.ResponseWriter, r *http.Request) {
	bucketName := r.FormValue("bucket")
	format := r.FormValue("format")
	fileName := r.FormValue("filename")

	// Check information sent are compliants
	if bucketName == "" {
		http.Error(w, "[400]- Bad request: bucket_name is missing", 400)
		return
	}
	if format == "" {
		http.Error(w, "[400]- Bad request: format is missing", 400)
		return
	}
	formatInt, _ := strconv.Atoi(format)
	err := encodeInMinio(bucketName, formatInt, fileName)
	if err != nil {
		http.Error(w, "[500]- Bad request: format is missing", 500)
		return
	}
	/* Send response */
	js, _ := json.Marshal(response{"Video successfully encoded"})

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

func encodeInMinio(bucketName string, maxQual int, fileName string) error {
	quals := []int{240, 360, 480, 720, 1080}
	cmpt := 0
	log.Printf("encode in minio")
	fileLoc, err := getFileInBucket(bucketName)
	if err != nil {
		return err
	}
	for quals[cmpt] < maxQual {
		// Encode to the format specified by quals
		target := fmt.Sprintf("%d_%s", quals[cmpt], fileName)
		cmd := fmt.Sprintf("ffmpeg -i %s -vf scale=-1:%d -c:v libx264 -crf 18 -preset veryslow -c:a copy %s", fileLoc, quals[cmpt], target)
		_, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			return err
		}

		file, err := os.Open(target) // For read access.
		if err != nil {
			log.Fatal(err)
		}

		// Save in minio
		info, err := (*config.Api.Minio).PutObject(context.Background(), bucketName, target, file, -1 /*fileStat.Size()*/, minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("Saved object: %s", info)
		cmpt++
	}
	return nil
}

func getFileInBucket(bucketName string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objectCh := (*config.Api.Minio).ListObjects(ctx, "mybucket", minio.ListObjectsOptions{
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			log.Println(object.Err)
			return "", object.Err
		}
		log.Println(object)
		break
	}
	// log.Println(objectCh)
	// log.Println(file)
	// err := (*config.Api.Minio).FGetObject(context.Background(), bucketName, file, "/tmp/tmp_file.mp4", minio.GetObjectOptions{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	return "", nil
}
