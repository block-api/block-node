package transporter

type Channel string

const (
	ChanDiscovery = "discovery"
	ChanMessage   = "msg"
)

type Transporter interface {
	Connect() error
	Disconnect() error
	Send(target string, payload interface{}) error
	Subscribe(channel Channel) error
}

type TransportPocket struct {
	target  string
	payload interface{}
}

func NewTransportPocket(target string, payload interface{}) TransportPocket {
	return TransportPocket{
		target,
		payload,
	}
}
