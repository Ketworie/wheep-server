package storage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func HandleUploadImage(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	// 10 MB
	resourceAddress, err := UploadImage(userId, r, 10)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\"" + resourceAddress + "\""))
	return err
}

func UploadImage(userId primitive.ObjectID, r *http.Request, maxMemory int64) (string, error) {
	err := r.ParseMultipartForm(maxMemory << 20)
	if err != nil {
		return "", err
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		return "", err
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Print(closeErr)
		}
	}()
	resourceAddress, err := Upload(userId, file)
	return resourceAddress, err
}
