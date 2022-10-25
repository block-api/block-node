package utils

type (
	NodeID    string
	NodeName  string
	BlockName string
)

func (name NodeID) String() string {
	return string(name)
}

func (name NodeName) String() string {
	return string(name)
}

func (name BlockName) String() string {
	return string(name)
}
