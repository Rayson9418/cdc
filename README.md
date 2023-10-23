# cdc

A go library to help you monitor changed data and sync data in Mysql/Mongo.

## how to use
With Go module support, simply add the following import

```
import "github.com/Rayson9418/cdc"
```

## demo

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