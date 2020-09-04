package slice

import "go.mongodb.org/mongo-driver/bson/primitive"

func HasString(ss []string, s string) bool {
	for _, entry := range ss {
		if s == entry {
			return true
		}
	}
	return false
}

func HasId(ids []primitive.ObjectID, id primitive.ObjectID) bool {
	for _, entry := range ids {
		if id == entry {
			return true
		}
	}
	return false
}
