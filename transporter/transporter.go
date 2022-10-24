package transporter

type Channel string

const (
	ChanPrefix    = "bn_"
	ChanDiscovery = ChanPrefix + "discovery"
	ChanMessage   = ChanPrefix + "msg"
	ChanStorage   = ChanPrefix + "store"
)

type Transporter interface {
	Connect() error
	Disconnect() error
	Send(channel Channel, payload []byte) error
	Subscribe(channel Channel, callback func(pocket Pocket[[]byte])) error
}

type PayloadMessage []struct {
	Data interface{} `json:"data"`
}

type Event uint

const (
	EventConnected Event = iota + 1
	EventDisconnected
	EventPing
	EventPong
)
