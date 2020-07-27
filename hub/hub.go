package hub

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Model struct {
	ID           primitive.ObjectID   `bson:"_id"`
	Name         string               `bson:"name"`
	Image        string               `bson:"image"`
	LastUpdate   time.Time            `bson:"lastUpdate"`
	Users        []primitive.ObjectID `bson:"users"`
	LastModified time.Time            `bson:"lastModified"`
}

func (h Model) View() View {
	return View{
		ID:           h.ID,
		Name:         h.Name,
		Image:        h.Image,
		LastUpdate:   h.LastUpdate,
		UserCount:    len(h.Users),
		LastModified: h.LastModified,
	}
}

type View struct {
	ID           primitive.ObjectID `json:"id"`
	Name         string             `json:"name"`
	Image        string             `json:"image"`
	LastUpdate   time.Time          `json:"lastUpdate"`
	UserCount    int                `json:"userCount"`
	LastModified time.Time          `json:"lastModified"`
}

type AddView struct {
	Name  string               `json:"name"`
	Image string               `json:"image"`
	Users []primitive.ObjectID `json:"users"`
}
