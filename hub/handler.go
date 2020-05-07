package hub

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"wheep-server/user"
)

func HandleAdd(u user.Model, w http.ResponseWriter, r *http.Request) error {
	var h Model
	err := json.NewDecoder(r.Body).Decode(&h)
	if err != nil {
		return err
	}
	service := GetService()
	h.Users = map[uuid.UUID]bool{u.ID: true}
	add, err := service.Add(h)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(View{
		ID:        add.ID,
		Name:      add.Name,
		UserCount: 0,
	})
	if err != nil {
		log.Println(err)
	}
	return nil
}

func HandleFindByUser(u user.Model, w http.ResponseWriter, r *http.Request) error {
	s := r.FormValue("id")
	userId, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	service := GetService()
	hubs, err := service.FindByUser(userId)
	if err != nil {
		return err
	}
	views := make([]View, len(hubs))
	for i, v := range hubs {
		views[i] = View{
			ID:        v.ID,
			Name:      v.Name,
			UserCount: len(v.Users),
		}
	}
	return json.NewEncoder(w).Encode(views)
}
