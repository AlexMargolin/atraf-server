package bucket

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const (
	UploadsBaseDir = "uploads"
	FSBucketType   = "filesystem"
)

type FSBucket struct{}

func (FSBucket) Type() string {
	return FSBucketType
}

func (FSBucket) SaveFile(name string, path string, file multipart.File) (string, error) {
	dir := fmt.Sprintf("%s/%s", UploadsBaseDir, path)
	filename := dir + name

	if _, err := os.Stat(dir); err != nil {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", err
		}
	}

	dst, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	if err = dst.Sync(); err != nil {
		return "", err
	}

	return filename, nil
}

func (FSBucket) PrependBucketURL(filename string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("BUCKET_URL"), filename)
}

func (FSBucket) ServeFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dir := http.Dir(UploadsBaseDir)
		handler := http.StripPrefix("/"+UploadsBaseDir, http.FileServer(dir))

		handler.ServeHTTP(w, r)
	}
}

func NewFSBucket() *FSBucket {
	return &FSBucket{}
}
