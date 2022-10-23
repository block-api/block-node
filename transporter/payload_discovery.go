package transporter

type PayloadDiscovery struct {
	NodeID  string              `json:"node_id"`
	Name    string              `json:"name"`
	Version uint                `json:"version"`
	Blocks  map[string][]string `json:"blocks"`
	Event   Event               `json:"event"`
}
