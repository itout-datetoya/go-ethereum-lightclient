package util

import (
	"encoding/hex"
	"errors"
	"strings"
	"strconv"
)

func HexstrTo32Bytes(hexString string) ([32]byte) {
	hexString = strings.TrimPrefix(hexString, "0x")

	byteArray, err := hex.DecodeString(hexString)
	if err != nil {
		panic(err)
	}

	if len(byteArray) != 32 {
		panic(errors.New("hexString is not 32 bytes"))
	} else {
		hash := [32]byte{}
		copy(hash[:], byteArray)
		return hash
	}
}

func HexstrToUint64(hexString string) (uint64) {
	hex, _ := strconv.ParseInt(hexString, 0, 64)
	return uint64(hex)
}

func GetBitFromBytes(b []byte, i uint64) bool {
	return (b[i>>3]>>(i&7))&1 == 1
}

func SetBitToBytes(b []byte, i uint64, v bool) {
	if bit := byte(1) << (i & 7); v {
		b[i>>3] |= bit
	} else {
		b[i>>3] &^= bit
	}
}
