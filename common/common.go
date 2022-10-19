package common

const (
	CmdExit Cmd = iota
	CmdStop
	CmdStart
)

type Cmd int64
type Address string
type Data interface{}
type DataBytes []byte
type Hash []byte
type Sig []byte
type Timestamp int64

func RemoveFromSlice(slice []interface{}, index int) []interface{} {
	sliceLen := len(slice)
	sliceLastIndex := sliceLen - 1

	if index != sliceLastIndex {
		slice[index] = slice[sliceLastIndex]
	}

	return slice[:sliceLastIndex]
}
