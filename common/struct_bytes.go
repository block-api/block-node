package common

import "fmt"

type StructBytes struct {
}

func (sb *StructBytes) Bytes() []byte {
	bytes := []byte{}
	fmt.Println(sb)
	return bytes
}
