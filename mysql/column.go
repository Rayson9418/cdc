package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-mysql-org/go-mysql/schema"
)

var KEmptyTime = time.Date(1970, 1, 1, 8, 0, 0, 0, time.Local)

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
