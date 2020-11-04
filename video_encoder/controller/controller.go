package controller

import (
	"config"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/PaulBarrie/video_encoder/config"
	minio "github.com/minio/minio-go/v7"
)

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
	// js, _ := json.Marshal(utils.RespVideo{"ok", video})

	w.WriteHeader(http.StatusCreated)
}

func encodeInMinio(bucketName string, maxQual int, fileName string) error {
	quals := []int{240, 360, 480, 720, 1080}
	cmpt := 0
	log.Printf("encode in minio")
	err := getFileInBucket(bucketName)
	if err != nil {
		return err
	}
	for quals[cmpt] < maxQual {
		cmpt++
	}
	return nil
}

func getFileInBucket(bucketName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objectCh := (*config.Api.Minio).ListObjects(ctx, "mybucket", minio.ListObjectOptions{
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			log.Println(object.Err)
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
	return nil
}
