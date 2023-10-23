package mongo

import "go.mongodb.org/mongo-driver/bson"

type StreamObject struct {
	Id                *WatchId `bson:"_id"`
	OperationType     string
	FullDocument      bson.Raw
	Ns                NS
	UpdateDescription *UpdateDescription
	DocumentKey       map[string]interface{}
}

// NS Changed database
type NS struct {
	Database   string `bson:"db"`
	Collection string `bson:"coll"`
}

// UpdateDescription 修改的document字段和值
type UpdateDescription struct {
	RemoveFields []string               `bson:"removeFields"`
	UpdateFields map[string]interface{} `bson:"updateFields"`
}

// WatchId for resume token
// Specifies the logical starting point for the new change stream
type WatchId struct {
	Data string `bson:"_data"`
}
