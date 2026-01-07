package config

import (
	"log"
	"os"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() *minio.Client {
   	endpoint := os.Getenv("MINIO_ENDPOINT")
   	accessKeyID := os.Getenv("MINIO_ACCESS_KEY_ID")
   	secretAccessKey := os.Getenv("MINIO_SECRET_ACCESS_KEY")
   	useSSL := false // false untuk development lokal
   
   	minioClient, err := minio.New(endpoint, &minio.Options{
   		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
   		Secure: useSSL,
   	})
   	if err != nil {
   		log.Fatalln("Gagal terhubung ke MinIO:", err)
   	}
   
   	log.Println("Koneksi MinIO berhasil.")
   	return minioClient
   }
