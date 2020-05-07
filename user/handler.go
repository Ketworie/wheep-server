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
