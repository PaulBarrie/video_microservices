package config

import (
	"fmt"
	"log"
	"os"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func (a *App) ConnectMinio() {
	endpoint := os.Getenv("MINIO_HOST")
	accessKeyID := os.Getenv("MINIO_ROOT_USER")
	secretAccessKey := os.Getenv("MINIO_ROOT_PASSWORD")
	useSSL := false
	endpoint = fmt.Sprintf("%s:%s", endpoint, "9000")
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
