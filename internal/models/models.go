// Package models provides the data models for the application.
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseModel struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type User struct {
	BaseModel `bson:",inline"`
	Name      string `json:"name" bson:"name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"-" bson:"password"`
}
