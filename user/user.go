package user

import "github.com/google/uuid"

type Model struct {
	ID       uuid.UUID `bson:"_id",json:"id"`
	Alias    string    `bson:"alias",json:"alias"`
	Login    string    `bson:"login",json:"login"`
	Password string    `bson:"password",json:"password"`
	Name     string    `bson:"name",json:"name"`
}

type View struct {
	Alias string `json:"alias"`
	Name  string `json:"name"`
}
