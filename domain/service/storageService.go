package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
)

type StorageService interface {
	Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
	GetPresignedURL(objectName string) (string, error)
	BuildPublicURL(objectName string) (string, error)
}

type minioStorage struct {
	client         *minio.Client
	bucket         string
	publicEndpoint string
}

// GetPresignedURL implements [StorageService].
func (m *minioStorage) GetPresignedURL(objectName string) (string, error) {
	expiry := time.Second * 600
	reqParams := make(url.Values)

	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(objectName)))

	presignedURL, err := m.client.PresignedGetObject(context.Background(), m.bucket, objectName, expiry, reqParams)
	if err != nil {
		log.Println("Gagal membuat Presigned URL:", err)
		return "", err
	}

	return presignedURL.String(), nil
}

func (m *minioStorage) BuildPublicURL(objectName string) (string, error) {
	if m.publicEndpoint == "" {
		return "", fmt.Errorf("MINIO_PUBLIC_ENDPOINT not set")
	}
	pub, err := url.Parse(m.publicEndpoint)
	if err != nil {
		return "", err
	}
	pub.Path = strings.TrimRight(pub.Path, "/") + "/" + strings.TrimPrefix(objectName, "/")
	return pub.String(), nil
}

func NewFileService(minioClient *minio.Client) StorageService {
	bucketName := os.Getenv("MINIO_BUCKET")
	if bucketName == "" {
		bucketName = "rentiva"
	}
	publicEndpoint := os.Getenv("MINIO_PUBLIC_ENDPOINT") // optional: e.g. https://minio.xxx.com

	ctx := context.Background()
	region := "us-east-1"

	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	if errBucketExists != nil {
		log.Printf("minio: warning cek bucket %s: %v", bucketName, errBucketExists)
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region}); err != nil {
			errResp := minio.ToErrorResponse(err)
			if errResp.Code != "BucketAlreadyOwnedByYou" && errResp.Code != "BucketAlreadyExists" {
				log.Printf("minio: gagal membuat bucket %s: %v", bucketName, err)
			}
		}
	}

	return &minioStorage{
		client:         minioClient,
		bucket:         bucketName,
		publicEndpoint: publicEndpoint,
	}
}

func (m *minioStorage) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Pastikan folder tersusun rapi
	dir := filepath.Dir(objectName)
	if dir == "." {
		objectName = fmt.Sprintf("uploads/%s", objectName)
	}

	_, err := m.client.PutObject(ctx, m.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	// hanya kembalikan object key
	return objectName, nil
}

// Helper untuk generate nama file unik
func UniqueObjectName(prefix, filename string) string {
	now := time.Now().UnixNano()
	cleanName := strings.ReplaceAll(filename, " ", "_")
	return fmt.Sprintf("%s/%d_%s", strings.TrimSuffix(prefix, "/"), now, cleanName)
}
