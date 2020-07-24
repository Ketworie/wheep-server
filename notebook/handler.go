package notebook

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"wheep-server/user"
)

func HandleGetContacts(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	contacts, err := GetService().GetContacts(userId)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(contacts)
}

func HandleAddContact(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	alias := r.FormValue("alias")
	if len(alias) > 0 {
		return errors.New("no alias specified")
	}
	model, err := user.GetService().GetByAlias(alias)
	if err != nil {
		return err
	}
	return GetService().AddContact(userId, model.ID)
}

func HandleRemoveContact(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	return GetService().AddContact(userId, id)
}
