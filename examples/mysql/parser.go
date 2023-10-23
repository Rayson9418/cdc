package mysql

import (
	"time"

	cdcmysql "github.com/Rayson9418/cdc/mysql"
)

type DemoData struct {
	Id        int
	EventTime time.Time
	EventName string
	EventDesc string
}

func ParseDemo(o *cdcmysql.RowEventData) *DemoData {
	return &DemoData{
		Id:        cdcmysql.GetInt(o.Row, "id"),
		EventTime: cdcmysql.GetTime(o.Row, "event_time"),
		EventName: cdcmysql.GetString(o.Row, "event_name"),
		EventDesc: cdcmysql.GetString(o.Row, "event_desc"),
	}
}
