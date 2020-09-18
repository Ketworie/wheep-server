package notebook

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func HandleGetContacts(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	contacts, err := GetRepository().GetContacts(userId)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	return json.NewEncoder(w).Encode(contacts)
}

func HandleAddContact(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	contactId, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	if userId == contactId {
		return errors.New("you cannot add yourself to contacts")
	}
	return GetRepository().AddContact(userId, contactId)
}

func HandleRemoveContact(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	return GetRepository().AddContact(userId, id)
}
