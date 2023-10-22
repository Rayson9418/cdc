package mysql

type RowEventData struct {
	Database  string
	Table     string
	Action    string
	Timestamp uint32
	Pos       uint32
	Old       map[string]interface{}
	Row       map[string]interface{}
	Err       error
}
