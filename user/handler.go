package user

import (
	"encoding/json"
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
	err = json.NewEncoder(w).Encode(View{
		Alias: u.Alias,
		Name:  u.Name,
	})
	return err
}
