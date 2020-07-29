package user

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
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

func HandleGet(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	u, err := GetService().Get(id)
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
