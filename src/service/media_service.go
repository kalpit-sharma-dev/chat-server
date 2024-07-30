package service

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kalpit-sharma-dev/chat-service/src/models"
	"github.com/kalpit-sharma-dev/chat-service/src/repository"
)

type MediaService struct {
	mediaRepo *repository.MediaRepository
	s3Client  *s3.S3
	bucket    string
}

func NewMediaService(mediaRepo *repository.MediaRepository, s3Client *s3.S3, bucket string) *MediaService {
	return &MediaService{
		mediaRepo: mediaRepo,
		s3Client:  s3Client,
		bucket:    bucket,
	}
}

func (service *MediaService) UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	defer file.Close()
	buffer := make([]byte, header.Size)
	file.Read(buffer)

	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(header.Filename))
	_, err := service.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(service.bucket),
		Key:                aws.String(fileName),
		ACL:                aws.String("public-read"),
		Body:               bytes.NewReader(buffer),
		ContentLength:      aws.Int64(header.Size),
		ContentType:        aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", service.bucket, fileName)
	return url, nil
}

func (service *MediaService) SaveMedia(media *models.Media) error {
	return service.mediaRepo.SaveMedia(media)
}
