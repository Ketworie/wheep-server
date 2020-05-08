package hub

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"wheep-server/user"
)

func HandleAdd(u user.Model, w http.ResponseWriter, r *http.Request) error {
	var av AddView
	err := json.NewDecoder(r.Body).Decode(&av)
	if err != nil {
		return err
	}
	service := GetService()
	userMap := make(map[primitive.ObjectID]bool)
	users := []primitive.ObjectID{u.ID}
	for _, v := range av.Users {
		if _, ok := userMap[v]; !ok {
			userMap[v] = true
			users = append(users, v)
		}
	}
	add, err := service.Add(Model{
		Name:  av.Name,
		Users: users,
	})
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(add.View())
	if err != nil {
		log.Println(err)
	}
	return nil
}

func HandleGet(u user.Model, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	model, err := GetService().Get(id)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(model.View())
	if err != nil {
		log.Println(err)
	}
	return nil
}

func HandleDelete(u user.Model, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	return GetService().Delete(id)
}

func HandleRename(u user.Model, w http.ResponseWriter, r *http.Request) error {
	var v View
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}
	service := GetService()
	err = service.Rename(Model{
		ID:   v.ID,
		Name: v.Name,
	})
	if err != nil {
		return err
	}
	get, err := service.Get(v.ID)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(get.View())
	return err
}

func HandleAddUsers(u user.Model, w http.ResponseWriter, r *http.Request) error {
	var users []primitive.ObjectID
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	service := GetService()
	err = service.AddUsers(id, users)
	if err != nil {
		return err
	}
	model, err := service.Get(id)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(model.View())
	if err != nil {
		log.Println(err)
	}
	return nil
}

func HandleRemoveUsers(u user.Model, w http.ResponseWriter, r *http.Request) error {
	var users []primitive.ObjectID
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	service := GetService()
	err = service.RemoveUsers(id, users)
	if err != nil {
		return err
	}
	model, err := service.Get(id)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(model.View())
	if err != nil {
		log.Println(err)
	}
	return nil
}

func HandleFindByUser(u user.Model, w http.ResponseWriter, r *http.Request) error {
	s := r.FormValue("id")
	userId, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return err
	}
	service := GetService()
	hubs, err := service.FindByUser(userId)
	if err != nil {
		return err
	}
	views := make([]View, len(hubs))
	for i, hub := range hubs {
		views[i] = hub.View()
	}
	return json.NewEncoder(w).Encode(views)
}

func HandleFindMyHubs(u user.Model, w http.ResponseWriter, r *http.Request) error {
	service := GetService()
	hubs, err := service.FindByUser(u.ID)
	if err != nil {
		return err
	}
	views := make([]View, len(hubs))
	for i, hub := range hubs {
		views[i] = hub.View()
	}
	return json.NewEncoder(w).Encode(views)
}
