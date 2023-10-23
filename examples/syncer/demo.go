package syncer

import (
	cdcsyncer "github.com/Rayson9418/cdc/syncer"

	"examples/options"
)

func DemoSyncAlways() error {
	opt := cdcsyncer.NewDefaultOpt()
	opt.StartHour = options.CdcOpt.Syncer.StartHour
	opt.EndHour = options.CdcOpt.Syncer.EndHour
	opt.BatchLimit = options.CdcOpt.Syncer.BatchLimit
	opt.Interval = options.CdcOpt.Syncer.Interval

	demo1Syncer := NewDemo1Syncer()

	return cdcsyncer.StartSyncer(demo1Syncer)
}

func DemoSyncOnce() error {
	opt := cdcsyncer.NewDefaultOpt()
	opt.StartHour = options.CdcOpt.Syncer.StartHour
	opt.EndHour = options.CdcOpt.Syncer.EndHour
	opt.BatchLimit = options.CdcOpt.Syncer.BatchLimit
	opt.Interval = options.CdcOpt.Syncer.Interval

	demo1Syncer := NewDemo1Syncer()

	return cdcsyncer.SyncOnce(demo1Syncer)
}

func DemoSyncOnTime() error {
	opt := cdcsyncer.NewDefaultOpt()
	opt.StartHour = options.CdcOpt.Syncer.StartHour
	opt.EndHour = options.CdcOpt.Syncer.EndHour
	opt.BatchLimit = options.CdcOpt.Syncer.BatchLimit
	opt.Interval = options.CdcOpt.Syncer.Interval

	demo1Syncer := NewDemo1Syncer()

	return cdcsyncer.StartSyncerOnTime(demo1Syncer)
}
