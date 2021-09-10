package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Employee struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Department  string             `json:"department,omitempty" bson:"department,omitempty"`
	Level       string             `json:"level,omitempty" bson:"level,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
}
