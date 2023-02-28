package repository

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
)

const (
	projectID  = "tixologi-partner"
	bucketName = "partner-portal-storage"
	pathPrefix = "https://storage.cloud.google.com/partner-portal-storage/"
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var uploader *ClientUploader

func (c *ClientUploader) SetupBucket(scope string) {
	pathFolder := calculatePathFolder(scope)
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	uploader = &ClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
		uploadPath: pathFolder,
	}
}

func calculatePathFolder(scope string) string {
	switch scope {
	case "prod":
		return "users-prod/"
	case "test":
		return "users-test/"
	}
	return "users-dev/"
}

func HandleImageUpload(imageBase64 string, path string, ctx *gin.Context) (string, error) {
	// The actual image starts after the ","
	i := strings.Index(imageBase64, ",")
	if i < 0 {
		return "", errors.New("image must be in format: data:image/png;base64,encodedimage")
	}
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(imageBase64[i+1:]))
	suffix, err := calulateSuffix(imageBase64)
	if err != nil {
		return "", err
	}
	completePath := calculateCompletePath(path, suffix)
	err = uploader.NewUploadFile(dec, completePath)
	if err != nil {
		return "", err
	}
	prefix := pathPrefix + uploader.uploadPath

	return prefix + completePath, nil
}

func HandleFileUpload(fileBase64 string, path string, ctx *gin.Context) (string, error) {
	// The actual image starts after the ","
	i := strings.Index(fileBase64, ",")
	if i < 0 {
		return "", errors.New("files must be in format: data:application/pdf;base64,base64File")
	}
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(fileBase64[i+1:]))
	fileType := fileBase64[5:20]
	if fileType != "application/pdf" {
		return "", errors.New("file must be pdf")
	}
	completePath := calculateCompletePath(path, ".pdf")
	err := uploader.NewUploadFile(dec, completePath)
	if err != nil {
		return "", err
	}
	prefix := pathPrefix + uploader.uploadPath

	return prefix + completePath, nil
}
func calculateCompletePath(path string, suffix string) string {
	replacedPath := strings.Replace(path, "/", "", -1)
	completePath := fmt.Sprintf(replacedPath + suffix)
	return completePath

}

func calulateSuffix(imageBase64 string) (string, error) {
	switch imageBase64[5:15] {
	case "image/png;":
		return ".png", nil
	case "image/jpeg":
		return ".jpeg", nil
	}
	return "", errors.New("this is not a valid image")
}

func (c *ClientUploader) NewUploadFile(dec io.Reader, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, dec); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

// UploadFile uploads an object
func (c *ClientUploader) UploadFile(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
