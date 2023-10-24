package models

import (
	"context"
	"time"

	"go.uber.org/zap"

	. "examples/logger"
	"examples/mysql"
)

func QueryDemo1CountByTime(start, end time.Time) (int64, error) {
	db := mysql.GetDB(context.Background())

	var cnt int64
	err := db.Table(KTableDemo1).
		Where("event_time >= ?", start).
		Where("event_time < ?", end).
		Count(&cnt).
		Error

	if err != nil {
		Logger.Warn("get all demo1 err", zap.Error(err))
		return 0, err
	}
	return cnt, nil
}

func QueryDemo1DataByTime(offset, limit int, start, end time.Time) (interface{}, int64, error) {
	db := mysql.GetDB(context.Background())

	dataList := make([]*Demo1Model, 0)
	err := db.Table(KTableDemo1).
		Where("event_time >= ?", start).
		Where("event_time < ?", end).
		Offset(offset).
		Limit(limit).
		Scan(&dataList).
		Error
	if err != nil {
		Logger.Warn("get all demo1 err", zap.Error(err))
		return nil, 0, err
	}
	return dataList, int64(len(dataList)), nil
}
