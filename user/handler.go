package user

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"wheep-server/storage"
)

func HandleAdd(w http.ResponseWriter, r *http.Request) error {
	var u Model
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return err
	}
	service := GetService()
	_, err = service.Add(u)
	return err
}

func HandleGetByAlias(w http.ResponseWriter, r *http.Request) error {
	alias := r.FormValue("alias")
	service := GetService()
	u, err := service.GetByAlias(alias)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(u.View())
	return err
}

func HandleUpdateAvatar(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	// 5 MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}
	file, _, err := r.FormFile("image")
	if err != nil {
		return err
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			log.Print(closeErr)
		}
	}()
	resourceAddress, err := storage.Upload(userId, file)
	err = GetService().UpdateAvatar(userId, resourceAddress)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\"" + resourceAddress + "\""))
	return err
}
func HandleGet(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	alias := r.FormValue("alias")
	if err != nil && len(alias) == 0 {
		return err
	}
	service := GetService()
	var u Model
	if len(alias) == 0 {
		u, err = service.Get(id)
	} else {
		u, err = service.GetByAlias(alias)
	}
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(u.View())
	return err
}

func HandleGetList(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var is []primitive.ObjectID
	err := json.NewDecoder(r.Body).Decode(&is)
	if err != nil {
		return err
	}
	us, err := GetService().GetList(is)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(us.View())
	return err
}
