package mysql

import (
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"go.uber.org/zap"

	. "github.com/Rayson9418/cdc/logger"
	cdcstore "github.com/Rayson9418/cdc/store"
)

type BinlogHandler interface {
	GetDbName() string
	GetTableName() string
	GetActions() []string
}

type OnRowHandler interface {
	BinlogHandler
	OnRow(*RowEventData) error
}

type DummyBinLogHandler struct {
	DbName    string
	TableName string
	Actions   []string
}

func (h *DummyBinLogHandler) GetDbName() string {
	return h.DbName
}

func (h *DummyBinLogHandler) GetTableName() string {
	return h.TableName
}

func (h *DummyBinLogHandler) GetActions() []string {
	return h.Actions
}

func (h *DummyBinLogHandler) OnRow(*RowEventData) error {
	return nil
}

type DispatchHandler struct {
	gTidSet         mysql.GTIDSet
	databases       []*Database
	tableMap        map[string]*TableInfo
	tableActionSet  map[string]map[string]struct{}     // database.table -> action
	onRowHandlerMap map[string]map[string]OnRowHandler // database.table -> action -> handler

	cdcstore.MysqlPosInterface
	canal.DummyEventHandler
}

func NewDispatchHandler(dbs []*Database) *DispatchHandler {
	h := new(DispatchHandler)
	h.databases = dbs

	h.setActionSet()
	return h
}

func (m *DispatchHandler) setStore(store cdcstore.MysqlPosInterface) {
	m.MysqlPosInterface = store
}

func (m *DispatchHandler) setActionSet() {
	tableActionSet := make(map[string]map[string]struct{})
	for _, db := range m.databases {
		for _, t := range db.Tables {
			tableName := getTableKey(db.Name, t.Name)

			actionSet := make(map[string]struct{})
			for _, ac := range t.Actions {
				actionSet[ac] = struct{}{}
			}
			tableActionSet[tableName] = actionSet
		}
	}
	m.tableActionSet = tableActionSet
}

// Init the column name mapping of the data table.
func (m *DispatchHandler) initTableMapping(c *canal.Canal) error {
	tableMap := make(map[string]*TableInfo)
	for _, db := range m.databases {
		for _, t := range db.Tables {
			tableInfo, err := c.GetTable(db.Name, t.Name)
			if err != nil {
				Logger.Error("get table info err",
					zap.String("database", db.Name),
					zap.String("table", t.Name),
					zap.Error(err))
				return err
			}
			t.ColumnMapping, err = buildMapping(tableInfo)
			if err != nil {
				Logger.Error("build mapping err",
					zap.String("database", db.Name),
					zap.String("table", t.Name),
					zap.Error(err))
				return err
			}
			tableMap[tableInfo.String()] = t
		}
	}
	m.tableMap = tableMap
	return nil
}

func (m *DispatchHandler) addOnRowHandler(handlers ...OnRowHandler) error {
	handlerMap := make(map[string]map[string]OnRowHandler)

	for _, h := range handlers {
		tName := getTableKey(h.GetDbName(), h.GetTableName())
		actionSet, ok := m.tableActionSet[tName]
		if !ok {
			return fmt.Errorf("not support table, db: %s, table: %s", h.GetDbName(), h.GetTableName())
		}

		action2HandlerMap := make(map[string]OnRowHandler, 0)
		for _, ac := range h.GetActions() {
			if _, ok = actionSet[ac]; !ok {
				return fmt.Errorf("not support action, db: %s, table: %s, action: %s",
					h.GetDbName(), h.GetTableName(), ac)
			}
			action2HandlerMap[ac] = h
		}

		handlerMap[tName] = action2HandlerMap
	}
	m.onRowHandlerMap = handlerMap
	return nil
}

func (m *DispatchHandler) getHandler(e *canal.RowsEvent) (OnRowHandler, bool) {
	action2HandlerMap, ok := m.onRowHandlerMap[e.Table.String()]
	if !ok {
		Logger.Warn("empty onRowHandler",
			zap.String("table", e.Table.String()))
		return nil, false
	}
	handler, ok := action2HandlerMap[e.Action]
	if !ok {
		Logger.Warn("not support action",
			zap.String("table", e.Table.String()),
			zap.String("action", e.Action))
		return nil, false
	}
	return handler, true
}

func (m *DispatchHandler) GetGTidSet() mysql.GTIDSet {
	return m.gTidSet
}

func (m *DispatchHandler) SetGTidSet(g mysql.GTIDSet) {
	m.gTidSet = g
}

func (m *DispatchHandler) OnRow(e *canal.RowsEvent) error {
	Logger.Info("on row event",
		zap.String("e.Table", e.Table.String()),
		zap.String("e.Action", e.Action),
		zap.Uint32("e.Header.ServerID", e.Header.ServerID),
		zap.Uint32("e.Header.LogPos", e.Header.LogPos),
		zap.Any("e.Header.EventType", e.Header.EventType))

	l := new(RowEventData)
	l.Database = e.Table.Schema
	l.Table = e.Table.Name
	l.Action = e.Action
	l.Timestamp = e.Header.Timestamp
	l.Pos = e.Header.LogPos

	oldMap, rowMap := m.convert2Map(e)
	if len(rowMap) == 0 {
		return nil
	}

	l.Old = oldMap
	l.Row = rowMap

	h, ok := m.getHandler(e)
	if !ok {
		return nil
	}

	return h.OnRow(l)
}

func (m *DispatchHandler) OnPosSynced(eh *replication.EventHeader, p mysql.Position, g mysql.GTIDSet, force bool) error {
	m.SetGTidSet(g)

	err := m.Save(p.Name, p.Pos)
	if err != nil {
		Logger.Warn("binlog save err", zap.Error(err))
		return err
	}

	Logger.Info("on pos synced",
		zap.String("position", p.String()),
		zap.String("gTid", g.String()))
	return nil
}

func (m *DispatchHandler) convert2Map(e *canal.RowsEvent) (map[string]interface{}, map[string]interface{}) {
	var (
		oldData = make([]interface{}, 0)
		rowData = make([]interface{}, 0)
		old     = make(map[string]interface{})
		row     = make(map[string]interface{})
	)

	// The update event will have two rows, one for the old version and one for the new version.
	if len(e.Rows) == 2 {
		rowData = e.Rows[1]
		oldData = e.Rows[0]
	} else {
		rowData = e.Rows[0]
	}

	tableInfo, ok := m.tableMap[e.Table.String()]
	if !ok {
		Logger.Info("not support action, skip",
			zap.String("table", e.Table.String()),
			zap.String("action", e.Action))
		return old, row
	}

	if len(oldData) > 0 {
		old = make(map[string]interface{}, len(tableInfo.ColumnMapping))
		for columnName, padding := range tableInfo.ColumnMapping {
			old[columnName] = convertColumnData(oldData[padding.ColumnIndex], padding.ColumnMetadata)
		}
	}

	row = make(map[string]interface{})
	for columnName, padding := range tableInfo.ColumnMapping {
		row[columnName] = convertColumnData(rowData[padding.ColumnIndex], padding.ColumnMetadata)
	}

	return old, row
}
