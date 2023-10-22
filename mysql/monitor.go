package mysql

import (
	"fmt"

	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"go.uber.org/zap"

	. "github.com/Rayson9418/cdc/logger"
	cdcstore "github.com/Rayson9418/cdc/store"
)

type RowEventMonitor struct {
	c *canal.Canal
	h *DispatchHandler
}

func NewRowEventMonitor(opt *Options) (*RowEventMonitor, error) {
	m := new(RowEventMonitor)

	if err := m.initCanal(opt); err != nil {
		return nil, err
	}

	if err := m.initDispatchHandler(opt); err != nil {
		return nil, err
	}

	return m, nil
}

func (m *RowEventMonitor) SetStore(store cdcstore.MysqlPosInterface) {
	m.h.setStore(store)
}

func (m *RowEventMonitor) AddHandler(handlers ...OnRowHandler) error {
	return m.h.addOnRowHandler(handlers...)
}

func (m *RowEventMonitor) StartMonitor() error {
	// Set Dispatcher handler
	if m.h == nil {
		return fmt.Errorf("dispatch handler nil")
	}
	m.c.SetEventHandler(m.h)

	// Retrieve the position where the last monitoring ended
	_file, _pos, err := m.h.Pos()
	if nil != err {
		Logger.Error("get bin log position err", zap.Error(err))
		return err
	}
	// If it is the first time listening, start from the current gTid
	// If listening from gTid fails, start from the current position
	if _pos == 0 {
		gTidSet, err := m.c.GetMasterGTIDSet()
		if err != nil {
			Logger.Error("get master gTid err", zap.Error(err))
			return err
		}

		if err = m.c.StartFromGTID(gTidSet); err != nil {
			Logger.Warn("canal run from gTid err", zap.Error(err))
			masterPos, err := m.c.GetMasterPos()
			if err != nil {
				Logger.Error("get master pos err", zap.Error(err))
				return err
			}
			if err = m.c.RunFrom(masterPos); err != nil {
				Logger.Warn("canal run from position err", zap.Error(err))
				return err
			}
			Logger.Info("run from pos after from gTid",
				zap.String("file", masterPos.Name),
				zap.Uint32("pos", masterPos.Pos))
			return nil
		}
		Logger.Info("run from gTid", zap.String("gTid", gTidSet.String()))
		return nil
	}

	pos := mysql.Position{
		Name: _file,
		Pos:  _pos,
	}

	// Start listening
	if err = m.c.RunFrom(pos); err != nil {
		Logger.Warn("canal run from position err", zap.Error(err))
		return err
	}

	Logger.Info("run from pos", zap.String("file", pos.Name),
		zap.Uint32("pos", pos.Pos))
	return nil
}

func (m *RowEventMonitor) initDispatchHandler(opt *Options) error {
	h := NewDispatchHandler(opt.Databases)

	if err := h.initTableMapping(m.c); err != nil {
		return err
	}

	m.h = h
	return nil
}

func (m *RowEventMonitor) initCanal(opt *Options) error {
	// Get the monitored database and table.
	tableRegexList := make([]string, 0, len(opt.Databases))
	databases := make([]string, 0, len(opt.Databases))
	for _, db := range opt.Databases {
		databases = append(databases, db.Name)
		for _, t := range db.Tables {
			tableRegexList = append(tableRegexList, getTableKey(db.Name, t.Name))
		}
	}

	// New canal
	cfg := canal.NewDefaultConfig()
	cfg.Addr = opt.Addr
	cfg.User = opt.User
	cfg.Password = opt.Pwd
	cfg.Dump.Databases = databases
	cfg.IncludeTableRegex = tableRegexList
	cfg.Dump.ExecutionPath = ""

	c, err := canal.NewCanal(cfg)
	if err != nil {
		Logger.Warn("new canal err", zap.Error(err))
		return err
	}

	m.c = c
	return nil
}
