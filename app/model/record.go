package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DataGroupRecord struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Account   string             `bson:"account,omitempty"`
	Value     string             `bson:"value,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
	ExpiredAt time.Time          `bson:"expired_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at,omitempty"`
}
