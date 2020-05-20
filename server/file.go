package server

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"wheep-server/user"
)

func Upload(u user.Model, file multipart.File) (string, error) {
	fileDir := u.ID.Hex()
	fileName := primitive.NewObjectID().Hex()
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
