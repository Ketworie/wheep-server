package security

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"wheep-server/user"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	gate := GetGate()
	login := r.PostFormValue("login")
	password := r.PostFormValue("password")
	u, err := gate.Login(login, password)
	if err != nil {
		return err
	}
	w.Header().Set("X-Auth-Token", u.Hex())
	return nil
}

func HandleCreateIndexes(uid primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	return user.GetService().CreateIndexes()
}

func HandleMe(uid primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	model, err := user.GetService().Get(uid)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(model.View())
}

func HandleAuthorize(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, error) {
	get := r.Header.Get("X-Auth-Token")
	if len(get) == 0 {
		return primitive.NilObjectID, errors.New("unauthorized")
	}
	id, err := primitive.ObjectIDFromHex(get)
	if err != nil {
		return primitive.NilObjectID, err
	}
	gate := GetGate()
	uid, err := gate.Authorize(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return uid, nil
}
