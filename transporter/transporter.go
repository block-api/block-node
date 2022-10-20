package transporter

type Transporter interface {
	Connect() error
	Disconnect() error
	Send(target string, payload interface{}) error
}

type Channel int

const (
	ChanDiscovery = iota + 1
	ChanMessage
)
