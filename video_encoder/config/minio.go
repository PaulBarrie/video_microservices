package config

import (
	"log"
	"os"
	"fmt"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (a *App) ConnectMinio() {
	endpoint := os.Getenv("MINIO_CONTAINER")
	accessKeyID := os.Getenv("MINIO_ACCESS")
	secretAccessKey := os.Getenv("MINIO_SECRET")
	useSSL := false
	endpoint = fmt.Sprintf("%s:%s", endpoint, "9000")
	log.Printf(endpoint)
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
        Secure: useSSL,
    })
	if err != nil {
		log.Fatalln(err)
	}
	a.Minio = minioClient
}
