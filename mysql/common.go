package mysql

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Rayson9418/cdc/common"
	. "github.com/Rayson9418/cdc/logger"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/schema"
)

func getTableKey(dbName, tName string) string {
	return fmt.Sprintf(common.KTableNameFmt, dbName, tName)
}

func convertColumnData(value interface{}, col *schema.TableColumn) interface{} {
	if value == nil {
		return nil
	}

	switch col.Type {
	case schema.TYPE_ENUM:
		switch value := value.(type) {
		case int64:
			eNum := value - 1
			if eNum < 0 || eNum >= int64(len(col.EnumValues)) {
				// we insert invalid enum value before, so return empty
				Logger.Warn(fmt.Sprintf("invalid binlog enum index %d, for enum %v", eNum, col.EnumValues))
				return ""
			}
			return col.EnumValues[eNum]
		case string:
			return value
		case []byte:
			return string(value)
		}
	case schema.TYPE_SET:
		switch value := value.(type) {
		case int64:
			bitmask := value
			sets := make([]string, 0, len(col.SetValues))
			for i, s := range col.SetValues {
				if bitmask&int64(1<<uint(i)) > 0 {
					sets = append(sets, s)
				}
			}
			return strings.Join(sets, ",")
		}
	case schema.TYPE_BIT:
		switch value := value.(type) {
		case string:
			if value == "\x01" {
				return int64(1)
			}
			return int64(0)
		}
	case schema.TYPE_STRING:
		switch value := value.(type) {
		case []byte:
			return string(value[:])
		}
	case schema.TYPE_JSON:
		var f interface{}
		var err error
		switch v := value.(type) {
		case string:
			err = json.Unmarshal([]byte(v), &f)
		case []byte:
			err = json.Unmarshal(v, &f)
		}
		if err == nil && f != nil {
			return f
		}
	case schema.TYPE_DATETIME, schema.TYPE_TIMESTAMP, schema.TYPE_DATE:
		var vv string
		switch v := value.(type) {
		case string:
			vv = v
		case []byte:
			vv = string(v)
		}

		vt, err := time.ParseInLocation(mysql.TimeFormat, vv, time.Local)
		if err != nil || vt.IsZero() { // failed to parse date or zero date
			return time.Date(1970, 1, 1, 8, 0, 0, 0, time.Local)
		}
		return vt
	case schema.TYPE_NUMBER:
		switch v := value.(type) {
		case string:
			vv, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				Logger.Error(err.Error())
				return nil
			}
			return vv
		case []byte:
			str := string(v)
			vv, err := strconv.ParseInt(str, 10, 64)
			if err != nil {
				Logger.Error(err.Error())
				return nil
			}
			return vv
		}
	case schema.TYPE_DECIMAL, schema.TYPE_FLOAT:
		switch v := value.(type) {
		case string:
			vv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				Logger.Error(err.Error())
				return nil
			}
			return vv
		case []byte:
			str := string(v)
			vv, err := strconv.ParseFloat(str, 64)
			if err != nil {
				Logger.Error(err.Error())
				return nil
			}
			return vv
		}
	}

	return value
}

func GetString(row map[string]interface{}, key string) string {
	v := row[key]
	if v == nil {
		return ""
	}

	switch v := v.(type) {
	case string:
		return v

	default:
		return fmt.Sprintf("%v", v)
	}
}

func GetInt(row map[string]interface{}, key string) int {
	v := row[key]
	if v == nil {
		return 0
	}

	switch v := v.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	default:
		return 0
	}
}

func GetTime(row map[string]interface{}, key string) time.Time {
	v := row[key]
	if v == nil {
		return KEmptyTime
	}

	switch v := v.(type) {
	case time.Time:
		return v

	default:
		return KEmptyTime
	}
}
