package user

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"wheep-server/storage"
)

func HandleAdd(w http.ResponseWriter, r *http.Request) error {
	var u Model
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		return err
	}
	repository := GetRepository()
	_, err = repository.Add(r.Context(), u)
	return err
}

func HandleGetByAlias(w http.ResponseWriter, r *http.Request) error {
	alias := r.FormValue("alias")
	repository := GetRepository()
	u, err := repository.GetByAlias(r.Context(), alias)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(u.View())
	return err
}

func HandleUpdateAvatar(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	// 5 MB
	resourceAddress, err := storage.UploadImage(userId, r, 5)
	if err != nil {
		return err
	}
	err = GetRepository().UpdateAvatar(r.Context(), userId, resourceAddress)
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
	repository := GetRepository()
	var u Model
	if len(alias) == 0 {
		u, err = repository.Get(r.Context(), id)
	} else {
		u, err = repository.GetByAlias(r.Context(), alias)
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
	us, err := GetRepository().GetList(r.Context(), is)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(us.View())
	return err
}
