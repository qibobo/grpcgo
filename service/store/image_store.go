package store

import (
	"bytes"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const image_path = "images"

type ImageStore interface {
	Save(id string, imageType string, data bytes.Buffer) (string, error)
}

type DiskImageStore struct {
}

func (dms *DiskImageStore) Save(id string, imageType string, data bytes.Buffer) (string, error) {
	uuid := uuid.New().String()
	imageFileName := filepath.Join(image_path, uuid+"."+imageType)
	imageFile, err := os.Create(imageFileName)
	if err != nil {
		log.Printf("failed to create image file for %s with err %s", imageFileName, err)
		return "", err
	}
	_, err = data.WriteTo(imageFile)
	if err != nil {
		log.Printf("failed to write to image file for %s with err %s", imageFileName, err)
		return "", err
	}
	return imageFileName, nil
}
