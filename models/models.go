package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Create Struct
type Member struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Lastname  string             `json:"lastname" bson:"lastname,omitempty"`
	Age    string         `json:"age" bson:"age,omitempty"`
	Job    string         `json:"job" bson:"job,omitempty"`
	Status    string         `json:"status" bson:"status,omitempty"`
	Address    string         `json:"address" bson:"address,omitempty"`
}

