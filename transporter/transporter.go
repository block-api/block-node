package transporter

import "github.com/block-api/block-node/common/types"

type Channel string

const (
	ChanPrefix                  = "bn_"
	ChanDiscovery       Channel = ChanPrefix + "discovery"
	ChanMessage         Channel = ChanPrefix + "msg"
	ChanMessageResponse Channel = ChanPrefix + "msg_res"
	ChanStorage         Channel = ChanPrefix + "storage"
)

const (
	RandomDelivery types.DeliveryMethod = "random"
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
