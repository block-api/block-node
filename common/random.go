package common

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	math_rand "math/rand"
)

func GenerateRandomUInt64() uint64 {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand")
	}

	math_rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	return math_rand.Uint64()
}

func GetSmallestUInt64FromArray(randArr []uint64) uint64 {
	var smallest uint64 = randArr[0]

	for _, num := range randArr {
		if num < smallest {
			smallest = num
		}
	}
	return smallest
}
