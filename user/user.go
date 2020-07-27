package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Model struct {
	ID           primitive.ObjectID `bson:"_id"`
	Alias        string             `bson:"alias"`
	Login        string             `bson:"login"`
	Password     string             `bson:"password"`
	Name         string             `bson:"name"`
	Image        string             `bson:"image"`
	LastModified time.Time          `bson:"lastModified"`
}

func (u Model) View() View {
	return View{
		ID:           u.ID,
		Alias:        u.Alias,
		Name:         u.Name,
		Image:        u.Image,
		LastModified: u.LastModified,
	}
}

type View struct {
	ID           primitive.ObjectID `json:"id"`
	Alias        string             `json:"alias"`
	Name         string             `json:"name"`
	Image        string             `json:"image"`
	LastModified time.Time          `json:"lastModified"`
}
