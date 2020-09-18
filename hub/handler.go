package hub

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"wheep-server/storage"
)

func HandleAdd(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var av AddView
	err := json.NewDecoder(r.Body).Decode(&av)
	if err != nil {
		return err
	}
	repository := GetRepository()
	userMap := make(map[primitive.ObjectID]bool)
	users := []primitive.ObjectID{userId}
	for _, v := range av.Users {
		if _, ok := userMap[v]; !ok {
			userMap[v] = true
			users = append(users, v)
		}
	}
	add, err := repository.Add(Model{
		Name:  av.Name,
		Image: av.Image,
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

func HandleGet(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	model, err := GetRepository().Get(id)
	if err != nil {
		return err
	}
	err = json.NewEncoder(w).Encode(model.View())
	if err != nil {
		log.Println(err)
	}
	return nil
}

func HandleDelete(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	id, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	return GetRepository().Delete(id)
}

func HandleUpdateAvatar(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("hubId"))
	if err != nil {
		return err
	}
	err = GetRepository().AssertMember(hubId, userId)
	if err != nil {
		return err
	}
	// 5 MB
	resourceAddress, err := storage.UploadImage(userId, r, 5)
	if err != nil {
		return err
	}
	err = GetRepository().UpdateAvatar(hubId, resourceAddress)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\"" + resourceAddress + "\""))
	return err
}

func HandleRename(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("id"))
	if err != nil {
		return err
	}
	repository := GetRepository()
	err = repository.AssertMember(hubId, userId)
	if err != nil {
		return err
	}
	name := r.FormValue("name")
	err = repository.Rename(hubId, name)
	return err
}

func HandleAddUser(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var users []primitive.ObjectID
	id, err := primitive.ObjectIDFromHex(r.FormValue("hubId"))
	if err != nil {
		return err
	}
	for _, s := range r.Form["userId"] {
		userId, err := primitive.ObjectIDFromHex(s)
		if err != nil {
			return err
		}
		users = append(users, userId)
	}
	repository := GetRepository()
	err = repository.AssertMember(id, userId)
	if err != nil {
		return err
	}
	err = repository.AddUsers(id, users)
	return err
}

func HandleRemoveUser(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	removed, err := primitive.ObjectIDFromHex(r.FormValue("userId"))
	if err != nil {
		return err
	}
	id, err := primitive.ObjectIDFromHex(r.FormValue("hubId"))
	if err != nil {
		return err
	}
	repository := GetRepository()
	err = repository.AssertMember(id, userId)
	if err != nil {
		return err
	}
	err = repository.RemoveUser(id, removed)
	return err
}

func HandleFindByUser(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	s := r.FormValue("id")
	userId, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return err
	}
	repository := GetRepository()
	hubs, err := repository.FindByUser(userId)
	if err != nil {
		return err
	}
	views := make([]View, len(hubs))
	for i, hub := range hubs {
		views[i] = hub.View()
	}
	return json.NewEncoder(w).Encode(views)
}

func HandleFindMyHubs(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	repository := GetRepository()
	hubs, err := repository.FindByUser(userId)
	if err != nil {
		return err
	}
	views := make([]View, len(hubs))
	for i, hub := range hubs {
		views[i] = hub.View()
	}
	return json.NewEncoder(w).Encode(views)
}
