package storage

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
)

func HandleUpload(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	// 10 MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}
	file, _, err := r.FormFile("data")
	if err != nil {
		return err
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Print(closeErr)
		}
	}()
	resourceAddress, err := Upload(userId, file)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\"" + resourceAddress + "\""))
	return err
}
