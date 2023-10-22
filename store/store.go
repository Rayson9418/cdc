package store

type MysqlPosInterface interface {
	Pos() (string, uint32, error)
	Save(string, uint32) error
}

type MongoPosInterface interface {
	Pos() (string, error)
	Save(string) error
}

type SyncerPosInterface interface {
	Pos() (*SyncerPos, error)
	Save(interface{}, ...int64) error
}

type SyncerPos struct {
	SyncStartTimestamp int64
	SyncEndTimestamp   int64
	Pos                int64
	LastSyncTime       int64
}
