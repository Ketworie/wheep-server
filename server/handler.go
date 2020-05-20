package server

import (
	"log"
	"net/http"
	"wheep-server/user"
)

func HandleUpload(u user.Model, w http.ResponseWriter, r *http.Request) error {
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
	resourceAddress, err := Upload(u, file)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(resourceAddress))
	return err
}
