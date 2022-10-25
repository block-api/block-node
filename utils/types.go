package utils

type (
	NodeID          string
	NodeName        string
	NodeVersionName string
	BlockName       string
	ActionName      string
)

func (nodeID NodeID) String() string {
	return string(nodeID)
}

func (name NodeName) String() string {
	return string(name)
}

func (name NodeVersionName) String() string {
	return string(name)
}

func (name BlockName) String() string {
	return string(name)
}

func (name ActionName) String() string {
	return string(name)
}
