package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-mysql-org/go-mysql/schema"
)

var KEmptyTime = time.Date(1970, 1, 1, 8, 0, 0, 0, time.Local)

//
//func initTableMapping(c *canal.Canal) (map[string]map[string]*TableInfo, error) {
//	tableMap := make(map[string]map[string]*TableInfo)
//	for _, m := range opt.Monitors {
//		for _, t := range m.Tables {
//			tableInfo, err := c.GetTable(m.Schema, t.Name)
//			if err != nil {
//				Logger.Error("get table info err",
//					zap.String("schema", m.Schema),
//					zap.String("table", t.Name),
//					zap.Error(err))
//				return nil, err
//			}
//			t.ColumnMapping, err = buildMapping(tableInfo)
//			if err != nil {
//				Logger.Error("build mapping err",
//					zap.String("schema", m.Schema),
//					zap.String("table", t.Name),
//					zap.Error(err))
//				return nil, err
//			}
//			ac2TableInfoMap := make(map[string]*TableInfo, len(t.Actions))
//			for _, ac := range t.Actions {
//				ac2TableInfoMap[ac] = t
//			}
//			opt.MonitorTableMap[tableInfo.String()] = ac2TableInfoMap
//		}
//	}
//	return tableMap, nil
//}

func buildMapping(t *schema.Table) (map[string]*columnInfo, error) {
	mappings := make(map[string]*columnInfo)
	for _, column := range t.Columns {
		index := tableColumn(t, column.Name)
		if index == -1 {
			return nil, fmt.Errorf("column[%s] not exist", column.Name)
		}
		_column := column
		mappings[column.Name] = &columnInfo{
			ColumnIndex:    index,
			ColumnMetadata: &_column,
		}
	}
	return mappings, nil
}

func tableColumn(t *schema.Table, field string) int {
	for index, c := range t.Columns {
		if strings.EqualFold(c.Name, field) {
			return index
		}
	}
	return -1
}
