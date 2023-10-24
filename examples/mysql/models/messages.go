package models

import "time"

const (
	KTableDemo1 = "demo1"
	KTableDemo2 = "demo2"
)

type Demo1Model struct {
	Id        int       `gorm:"column:id"`
	EventTime time.Time `gorm:"column:event_time"`
	EventName string    `gorm:"column:event_name"`
	EventDesc string    `gorm:"column:event_desc"`
}

func (Demo1Model) TableName() string {
	return KTableDemo1
}

type Demo2Model struct {
	Id      int    `gorm:"column:id"`
	Name    string `gorm:"column:name"`
	Age     int    `gorm:"column:age"`
	Email   string `gorm:"column:email"`
	Address string `gorm:"column:address"`
}

func (Demo2Model) TableName() string {
	return KTableDemo2
}
