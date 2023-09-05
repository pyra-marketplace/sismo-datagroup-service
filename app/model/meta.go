package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DataGroupMate struct {
	Id                primitive.ObjectID `bson:"_id" json:"id"`
	GroupName         string             `bson:"group_name,omitempty"`
	Description       string             `bson:"description,omitempty"`
	Spec              string             `bson:"spec,omitempty"`
	TotalRecords      int                `bson:"total_records,omitempty"`
	StartAt           time.Time          `bson:"start_at,omitempty"`
	UpdatedAt         time.Time          `bson:"updated_at,omitempty"`
	GenerateFrequency string             `bson:"generate_frequency,omitempty"`
	Handler           string             `bson:"handler,omitempty"`
	TwitterConfig     TwitterConfig      `json:"config"`
}

type TwitterConfig struct {
	Followers int `json:"followers"`
}
