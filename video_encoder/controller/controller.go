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
	"strings"

	minio "github.com/minio/minio-go/v7"
)

type response struct {
	message string `json:"message"`
}

// EncodeVideo allows to encode video in Minio
func EncodeVideo(w http.ResponseWriter, r *http.Request) {
	bucketName := r.FormValue("bucket_name")
	format := r.FormValue("format")
	fileName := r.FormValue("filename")
	log.Printf("bucket: %s, format: %s, file: %s", bucketName, format, fileName)
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
		http.Error(w, "[500]- Internal server error", 500)
		return
	}
	/* Send response */
	js, _ := json.Marshal(response{"Video successfully encoded"})

	w.WriteHeader(http.StatusCreated)
	w.Write(js)
}

func encodeInMinio(bucketName string, maxQual int, fileName string) error {
	quals := []int{240, 360, 480, 720, 1080}
	fileBrut := strings.Split(fileName, "_")[1]
	cmpt := 0
	log.Printf("[+] Encoding in minio...")
	err := getFileInBucket(bucketName, fileName)
	if err != nil {
		log.Println(err)
		return err
	}

	for quals[cmpt] < maxQual {
		// Encode to the format specified by quals
		target := fmt.Sprintf("%d_%s", quals[cmpt], fileBrut)
		cmd := fmt.Sprintf("ffmpeg -i /tmp/tmp_file.mp4 -vf scale=%d:-2 -c:v libx264 -crf 18 -preset veryslow -c:a copy %s", quals[cmpt], target)
		out, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
		if err != nil {
			log.Println("Error in format change")
			log.Println(string(out))
			log.Println(err)
			return err
		}

		file, err := os.Open(target) // For read access.
		if err != nil {
			log.Println("Error in open target")
			log.Fatal(err)
		}

		// Save in minio
		info, err := (*config.API.Minio).PutObject(context.Background(), bucketName, target, file, -1 /*fileStat.Size()*/, minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if err != nil {
			log.Println("Error in put object")
			log.Println(err)
			return err
		}
		log.Printf("Saved object: %s", info)
		cmpt++
	}
	return nil
}

func getFileInBucket(bucketName string, fileName string) error {
	err := (*config.API.Minio).FGetObject(context.Background(), bucketName, fileName, "/tmp/tmp_file.mp4", minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Println("File saved !")
	return nil
}
