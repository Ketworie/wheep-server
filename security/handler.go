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

func HandleCreateIndexes(u user.Model, w http.ResponseWriter, r *http.Request) error {
	return user.GetService().CreateIndexes()
}

func HandleMe(u user.Model, w http.ResponseWriter, r *http.Request) error {
	return json.NewEncoder(w).Encode(u.View())
}

func HandleAuthorize(w http.ResponseWriter, r *http.Request) (user.Model, error) {
	get := r.Header.Get("X-Auth-Token")
	if len(get) == 0 {
		return user.Model{}, errors.New("unauthorized")
	}
	id, err := primitive.ObjectIDFromHex(get)
	if err != nil {
		return user.Model{}, err
	}
	gate := GetGate()
	u, err := gate.Authorize(id)
	if err != nil {
		return user.Model{}, err
	}
	return u, nil
}
