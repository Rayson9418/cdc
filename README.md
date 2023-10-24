# cdc

A go library to help you monitor changed data and sync data from Mysql/Mongo.

## how to use

With Go module support, simply add the following import

```
import "github.com/Rayson9418/cdc"
```

## demo

check the dir `examples` for more detail 

### monitor mongo changed data
```go
    func DemoMonitor() error {
        // Init position store
        if err := cdcredis.InitClient(options.CdcOpt.Redis); err != nil {
            Logger.Fatal("init redis client with opt err", zap.Error(err))
        }
        store := cdcstore.NewStreamRedisStore("mongo:oplog:pos")
    
        // New handler for specific collection
        handler := NewDemoHandler()
    
        // New mongo monitor
        m, err := cdcmongo.NewDefaultMonitor(options.CdcOpt.Mongo)
        if err != nil {
            return err
        }
        // Set position store for monitor
        m.SetStore(store)
        // Add handlers for monitor
        if err = m.AddHandler(handler); err != nil {
            return err
        }
        return m.StartMonitor()
    }
```

### monitor mysql changed data

```go
    func DemoMonitor() error {
        // InitClient position store
        if err := cdcredis.InitClient(options.CdcOpt.Redis); err != nil {
        Logger.Fatal("init redis client with opt err", zap.Error(err))
        }
        store := cdcstore.NewBinlogRedisStore("mysql:binlog:pos")
        
        // New handler for specific collection
        handler := NewDemoHandler()
        
        // New row event monitor
        m, err := cdcmysql.NewRowEventMonitor(options.CdcOpt.Mysql)
        if err != nil {
        return err
        }
        // Set position store for monitor
        m.SetStore(store)
        // Add handlers for monitor
        if err = m.AddHandler(handler); err != nil {
        return err
        }
        return m.StartMonitor()
    }
```

### sync data from mongo/mysql

```go
    // Keep sync always 
    func DemoSyncAlways() error {
        opt := cdcsyncer.NewDefaultOpt()
        opt.StartHour = options.CdcOpt.Syncer.StartHour
        opt.EndHour = options.CdcOpt.Syncer.EndHour
        opt.BatchLimit = options.CdcOpt.Syncer.BatchLimit
        opt.Interval = options.CdcOpt.Syncer.Interval
        
        demo1Syncer := NewDemo1Syncer()
        
        return cdcsyncer.StartSyncer(demo1Syncer)
    }
    
    // Sync once
    func DemoSyncOnce() error {
        opt := cdcsyncer.NewDefaultOpt()
        opt.StartHour = options.CdcOpt.Syncer.StartHour
        opt.EndHour = options.CdcOpt.Syncer.EndHour
        opt.BatchLimit = options.CdcOpt.Syncer.BatchLimit
        opt.Interval = options.CdcOpt.Syncer.Interval
        
        demo1Syncer := NewDemo1Syncer()
        
        return cdcsyncer.SyncOnce(demo1Syncer)
    }

    // Sync on time
    func DemoSyncOnTime() error {
        opt := cdcsyncer.NewDefaultOpt()
        opt.StartHour = options.CdcOpt.Syncer.StartHour
        opt.EndHour = options.CdcOpt.Syncer.EndHour
        opt.BatchLimit = options.CdcOpt.Syncer.BatchLimit
        opt.Interval = options.CdcOpt.Syncer.Interval
        
        demo1Syncer := NewDemo1Syncer()
        
        return cdcsyncer.StartSyncerOnTime(demo1Syncer)
    }
```

## run

If you have Docker and Docker Compose installed, you can run the code using the following steps.

1. Clone this repository to your machine, and switch to the directory `examples`.
    
   ```bash
   # clone this repository
   git clone git@github.com:Rayson9418/cdc.git
   
   # switch to the directory
   cd cdc/examples
   ```

2. Run the script to build the container for compilation, and set up the middleware environment.

   ```bash
   bash build_on_docker.sh -t
   ```
   
3. Compile the code and then run it.
   ```bash
   make build && ./examples
   ```