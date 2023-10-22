package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"time"
)

type DemoData struct {
	Id        primitive.ObjectID `json:"id"        bson:"_id"`
	EventTime time.Time          `json:"event_time" bson:"event_time"`
	EventName string             `json:"event_name" bson:"event_name"`
	EventDesc string             `json:"event_desc" bson:"event_desc"`
}

func ParseDemo(fullDoc bson.Raw) *DemoData {
	data := &DemoData{}

	if err := bson.Unmarshal(fullDoc, data); err != nil {

	}
	return data

}
