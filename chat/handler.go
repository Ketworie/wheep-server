package chat

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"wheep-server/hub"
	"wheep-server/message"
)

func HandleSend(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var vm message.View
	err := json.NewDecoder(r.Body).Decode(&vm)
	if err != nil {
		return err
	}
	vm.UserId = userId
	err = hub.GetRepository().AssertMember(vm.HubId, vm.UserId)
	if err != nil {
		return err
	}
	send, err := message.GetRepository().Add(message.Model{
		ID:     primitive.ObjectID{},
		UserId: vm.UserId,
		HubId:  vm.HubId,
		Text:   vm.Text,
		Date:   time.Now(),
		PrevId: primitive.ObjectID{},
	})
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(send.View())
}

func HandleLast(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("hub"))
	if err != nil {
		return err
	}
	last, err := message.GetRepository().Last(hubId)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(last.View())
}

func HandlePrev(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("hub"))
	if err != nil {
		return err
	}
	date, err := time.Parse("2006-01-02T15:04:05.999Z[MST]", r.FormValue("date"))
	if err != nil {
		return err
	}
	prev, err := message.GetRepository().Prev(hubId, date)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(prev.View())
}

func HandleNext(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("hub"))
	if err != nil {
		return err
	}
	date, err := time.Parse("2006-01-02T15:04:05.999Z[MST]", r.FormValue("date"))
	if err != nil {
		return err
	}
	last, err := message.GetRepository().Next(hubId, date)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(last.View())
}
