package mongo

type StreamHandler interface {
	GetDbName() string
	GetCollName() string
	GetOpTypes() []string
	OnChange(*StreamObject) error
}

type DummyStreamHandler struct {
	DbName   string
	CollName string
	OpTypes  []string
}

func (h *DummyStreamHandler) GetDbName() string {
	return h.DbName
}

func (h *DummyStreamHandler) GetCollName() string {
	return h.CollName
}

func (h *DummyStreamHandler) GetOpTypes() []string {
	return h.OpTypes
}

func (h *DummyStreamHandler) OnChange(*StreamObject) error {
	return nil
}
