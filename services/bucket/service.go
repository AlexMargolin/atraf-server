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
	SaveFile(name string, path string, file multipart.File) (string, error)
	PrependBucketURL(filename string) string
}

type Service struct {
	bucket Bucket
}

func (s Service) Save(f multipart.File) (string, error) {
	buffer := make([]byte, 512)
	if _, err := f.Read(buffer); err != nil {
		return "", err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)
	if !s.checkContentType(contentType) {
		return "", errors.New("unsupported content-type")
	}

	filename, filepath, err := s.uploadLocation(contentType)
	if err != nil {
		return "", err
	}

	path, err := s.bucket.SaveFile(filename, filepath, f)
	if err != nil {
		return "", err
	}

	return path, nil
}

func (s Service) FileURL(filename string) string {
	return s.bucket.PrependBucketURL(filename)
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
