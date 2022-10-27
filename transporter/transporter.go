package transporter

type Channel string

const (
	ChanPrefix    = "bn_"
	ChanDiscovery = ChanPrefix + "discovery"
	ChanMessage   = ChanPrefix + "msg"
	ChanStorage   = ChanPrefix + "storage"
)

type Transporter interface {
	Connect() error
	Disconnect() error
	Send(channel Channel, payload []byte) error
	Subscribe(channel Channel, callback func(payload []byte)) error
}

type Event uint

const (
	EventConnected Event = iota + 1
	EventDisconnected
	EventHeartbeat
)
