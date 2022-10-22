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
	Send(channel Channel, payload interface{}) error
	Subscribe(channel Channel) error
}

type TransportPocket struct {
	Target  string
	Payload interface{}
}

func NewTransportPocket(target string, payload interface{}) TransportPocket {
	return TransportPocket{
		Target:  target,
		Payload: payload,
	}
}
