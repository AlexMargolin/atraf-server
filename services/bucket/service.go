package bucket

import (
	"errors"
	"mime"
	"mime/multipart"
	"net/http"
	"time"

	"atraf-server/pkg/uid"
)

var allowedContentTypes = []string{
	"image/png",
	"image/jpeg",
}

type Bucket interface {
	Type() string
	SaveFile(name string, path string, file multipart.File) (string, error)
}

type Service struct {
	bucket Bucket
}

func (service Service) Type() string {
	return service.bucket.Type()
}

func (service Service) Save(file multipart.File) (string, error) {
	buffer := make([]byte, 512)
	if _, err := file.Read(buffer); err != nil {
		return "", err
	}

	// returns cursor back to start for any upcoming reads
	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	if !service.checkContentType(contentType) {
		return "", errors.New("unsupported content-type")
	}

	filename, filepath, err := service.uploadLocation(contentType)
	if err != nil {
		return "", err
	}

	path, err := service.bucket.SaveFile(filename, filepath, file)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (Service) uploadLocation(contentType string) (string, string, error) {
	extensions, err := mime.ExtensionsByType(contentType)
	if err != nil || len(extensions) == 0 {
		return "", "", err
	}

	filename := uid.New().String() + extensions[0]
	filepath := time.Now().Format("2006/01/02")

	return filename, filepath, nil
}

func (Service) checkContentType(contentType string) bool {
	for _, ct := range allowedContentTypes {
		if ct == contentType {
			return true
		}
	}

	return false
}

func NewService(b Bucket) *Service {
	return &Service{b}
}
