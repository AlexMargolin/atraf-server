package bucket

import (
	"fmt"
	"io"
	"mime/multipart"
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

func NewFSBucket() *FSBucket {
	return &FSBucket{}
}
