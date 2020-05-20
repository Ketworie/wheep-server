package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID       primitive.ObjectID `bson:"_id"`
	Alias    string             `bson:"alias"`
	Login    string             `bson:"login"`
	Password string             `bson:"password"`
	Name     string             `bson:"name"`
	Image    string             `bson:"image"`
}

func (u Model) View() View {
	return View{
		Alias: u.Alias,
		Name:  u.Name,
		Image: u.Image,
	}
}

type View struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
	Image string `json:"image"`
}
