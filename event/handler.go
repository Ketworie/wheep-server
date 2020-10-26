package event

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
	wheepTime "wheep-server/time"
)

func HandleLast(userId primitive.ObjectID, w http.ResponseWriter, r *http.Request) error {
	date, err := time.Parse(wheepTime.Zoned, r.FormValue("from"))
	if err != nil {
		return err
	}
	last, err := GetRepository().Last(r.Context(), userId, date)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(last.View())
}
