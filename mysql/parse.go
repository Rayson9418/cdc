package mysql

import "time"

type DemoData struct {
	Id         int
	ExpireTime time.Time
	Key        string
}

func ParseDemo(o *RowEventData) *DemoData {
	return &DemoData{
		Id:         GetInt(o.Row, "id"),
		ExpireTime: GetTime(o.Row, "expire_time"),
		Key:        GetString(o.Row, "key"),
	}
}
