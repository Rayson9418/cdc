package mongo

import (
	"context"
	"fmt"
	"time"

	cdcstore "github.com/Rayson9418/cdc/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	. "github.com/Rayson9418/cdc/logger"
)

type StreamMonitor struct {
	cdcstore.MongoPosInterface
	watcher     *mongo.ChangeStream
	nsOpTypeSet map[NS]map[string]struct{}
	handlerMap  map[NS]map[string]StreamHandler // database.collection -> insert/update/delete/ -> handler
}

func NewDefaultMonitor(opt *Options) (*StreamMonitor, error) {
	if err := InitClient(opt); err != nil {
		Logger.Error("init mongo client err", zap.Error(err))
		return nil, err
	}

	m := new(StreamMonitor)
	m.setActionSet(opt)
	return m, nil
}

func (m *StreamMonitor) setActionSet(opt *Options) {
	nsOpTypeSet := make(map[NS]map[string]struct{})
	for _, db := range opt.Databases {
		for _, coll := range db.Collections {
			ns := NS{
				Database:   db.Name,
				Collection: coll.Name,
			}
			opTypeSet := make(map[string]struct{})
			for _, ac := range coll.Actions {
				opTypeSet[ac] = struct{}{}
			}
			nsOpTypeSet[ns] = opTypeSet
		}
	}
	m.nsOpTypeSet = nsOpTypeSet
}

func (m *StreamMonitor) SetStore(store cdcstore.MongoPosInterface) {
	m.MongoPosInterface = store
}

func (m *StreamMonitor) AddHandler(handlers ...StreamHandler) error {
	handlerMap := make(map[NS]map[string]StreamHandler)
	for _, h := range handlers {
		ns := NS{
			Database:   h.GetDbName(),
			Collection: h.GetCollName(),
		}
		opTypeSet, ok := m.nsOpTypeSet[ns]
		if !ok {
			return fmt.Errorf("not support ns, db: %s, coll: %s", ns.Database, ns.Collection)
		}

		opType2HandlerMap := make(map[string]StreamHandler, 0)
		for _, opType := range h.GetOpTypes() {
			if _, ok = opTypeSet[opType]; !ok {
				return fmt.Errorf("not support opType, db: %s, coll: %s, op: %s", ns.Database, ns.Collection, opType)
			}
			opType2HandlerMap[opType] = h
		}

		handlerMap[ns] = opType2HandlerMap
	}
	m.handlerMap = handlerMap
	return nil
}

func (m *StreamMonitor) GetHandler(stream *StreamObject) (StreamHandler, bool) {
	action2HandlerMap, ok := m.handlerMap[stream.Ns]
	if !ok {
		Logger.Warn("not support collection", zap.Any("ns", stream.Ns))
		return nil, false
	}
	handler, ok := action2HandlerMap[stream.OperationType]
	if !ok {
		Logger.Warn("not support operation", zap.Any("ns", stream.Ns),
			zap.String("operation", stream.OperationType))
		return nil, false
	}
	return handler, true
}

func (m *StreamMonitor) SetWatcher() error {
	var resumeToken bson.Raw
	token, err := m.Pos()
	if err != nil {
		return err
	}
	if token != "" {
		Logger.Info("get watch token:", zap.String("token", token))
		resumeToken = []byte(fmt.Sprintf(tokenTmp, token))
	}

	// Set up the stream option.
	timestamp := &primitive.Timestamp{
		T: uint32(time.Now().Unix()),
		I: 0,
	}
	streamOpt := options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(timestamp)
	if len(resumeToken) != 0 {
		streamOpt.SetResumeAfter(resumeToken)
		streamOpt.SetStartAtOperationTime(nil)
	}

	// Set up stream filtering conditions.
	orCond := make([]bson.M, 0)
	for ns, actionSet := range m.nsOpTypeSet {
		optMatches := make([]string, 0, len(actionSet))
		for action := range actionSet {
			optMatches = append(optMatches, action)
		}

		orCond = append(orCond, bson.M{
			"ns": bson.M{
				"db":   ns.Database,
				"coll": ns.Collection,
			},
			"operationType": bson.M{"$in": optMatches},
		})
	}
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.M{"$or": orCond}}},
	}

	// Obtain the change stream.
	watcher, err := globalClient.Watch(context.Background(), pipeline, streamOpt)
	if err != nil {
		Logger.Error("watch coll failed", zap.Error(err))
		return err
	}

	m.watcher = watcher
	return nil
}

func (m *StreamMonitor) StartMonitor() error {
	if err := m.SetWatcher(); err != nil {
		return err
	}

	w := m.watcher
	for w.Next(context.Background()) {
		var stream StreamObject
		err := w.Decode(&stream)
		if err != nil {
			Logger.Error("watch Decode data failed",
				zap.String("current", w.Current.String()),
				zap.Error(err))
			continue
		}
		resumeToken := w.ResumeToken()

		Logger.Debug("=====> receive change data:",
			zap.Any("ns", stream.Ns),
			zap.Any("id", stream.DocumentKey),
			zap.Any("update", stream.UpdateDescription))

		handler, ok := m.GetHandler(&stream)
		if !ok {
			continue
		}

		if err = handler.OnChange(&stream); err != nil {
			return err
		}

		if err = m.Save(resumeToken.Lookup("_data").StringValue()); err != nil {
			return err
		}
	}
	return nil
}
