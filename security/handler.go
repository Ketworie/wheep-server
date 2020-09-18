package security

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"wheep-server/user"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	gate := GetGate()
	login := r.PostFormValue("login")
	password := r.PostFormValue("password")
	token, err := gate.Login(login, password)
	if err != nil {
		return err
	}
	w.Header().Set("X-Auth-Token", token)
	return nil
}

func HandleCreateIndexes(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	return user.GetRepository().CreateIndexes()
}

func HandleMe(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	model, err := user.GetRepository().Get(userId)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(model.View())
}

func HandleAuthorize(w http.ResponseWriter, r *http.Request) (primitive.ObjectID, error) {
	token := r.Header.Get("X-Auth-Token")
	if len(token) == 0 {
		return primitive.NilObjectID, errors.New("unauthorized")
	}
	sid, err := uuid.Parse(token)
	if err != nil {
		return primitive.NilObjectID, err
	}
	gate := GetGate()
	uid, err := gate.Authorize(sid)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return uid, nil
}
