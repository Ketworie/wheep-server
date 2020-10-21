package chat

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	"wheep-server/event"
	"wheep-server/hub"
	"wheep-server/message"
	wheepTime "wheep-server/time"
)

func HandleSend(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	var vm message.View
	err := json.NewDecoder(r.Body).Decode(&vm)
	if err != nil {
		return err
	}
	vm.UserId = userId
	err = hub.GetRepository().AssertMember(r.Context(), vm.HubId, vm.UserId)
	if err != nil {
		return err
	}
	model := message.Model{
		ID:     primitive.ObjectID{},
		UserId: vm.UserId,
		HubId:  vm.HubId,
		Text:   vm.Text,
		Date:   time.Now(),
		PrevId: primitive.ObjectID{},
	}
	send, err := message.GetRepository().Add(r.Context(), model)
	if err != nil {
		return err
	}
	view := send.View()
	eventModel, err := event.NewEventModel(view, view.UserId)
	if err != nil {
		return err
	}
	GetService().Fanout(r.Context(), view.HubId, eventModel.View())
	return json.NewEncoder(w).Encode(view)
}

func HandleSetup(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	token := r.Header.Get("X-Auth-Token")
	return GetService().SetupExchange(userId, token)
}

func HandlePrev(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	hubId, err := primitive.ObjectIDFromHex(r.FormValue("hub"))
	if err != nil {
		return err
	}
	date, err := time.Parse(wheepTime.Zoned, r.FormValue("date"))
	if err != nil {
		return err
	}
	prev, err := message.GetRepository().Prev(r.Context(), hubId, date)
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
	date, err := time.Parse(wheepTime.Zoned, r.FormValue("date"))
	if err != nil {
		return err
	}
	last, err := message.GetRepository().Next(r.Context(), hubId, date)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(last.View())
}
