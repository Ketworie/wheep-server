package storage

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
)

var ResourceRoot = "/resources/"

func Upload(userId primitive.ObjectID, file multipart.File) (string, error) {
	fileDir := userId.Hex()
	fileName := uuid.New().String()
	fileExtension := ".jpg"
	filePath := path.Join(ResourceRoot, fileDir, fileName+fileExtension)
	err := os.MkdirAll(path.Join(ResourceRoot, fileDir), os.ModePerm)
	if err != nil {
		return "", err
	}
	all, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filePath, all, os.ModePerm)
	if err != nil {
		return "", err
	}
	return fileDir + "/" + fileName + fileExtension, err
}
