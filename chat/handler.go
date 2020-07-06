package chat

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"wheep-server/message"
)

func HandleSend(uid primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var vm message.View
	err := json.NewDecoder(r.Body).Decode(&vm)
	if err != nil {
		return err
	}
	vm.UserId = uid
	return Send(vm)
}

func HandleLast(uid primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("hub"))
	if err != nil {
		return err
	}
	last, err := message.GetService().Last(hubId)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(last)
}

func HandlePrev(uid primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("hub"))
	if err != nil {
		return err
	}
	date, err := time.Parse(time.RFC3339, r.FormValue("date"))
	last, err := message.GetService().Prev(hubId, date)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(last)
}

func HandleNext(uid primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("hub"))
	if err != nil {
		return err
	}
	date, err := time.Parse(time.RFC3339, r.FormValue("date"))
	last, err := message.GetService().Next(hubId, date)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(last)
}
