package beacon

import (
	"encoding/hex"
	"errors"
	"strings"
	"strconv"
	"github.com/tidwall/gjson"
)

type BeaconBlock struct {
	slot uint64
	hash [32]byte
	parentHash [32]byte
	stateRoot [32]byte
	bodyRoot [32]byte
}

func Parse(data string) (BeaconBlock) {
	block := BeaconBlock{}

	block.slot = hexstrToUint64(gjson.Get(data, "data.header.message.slot").String())
	block.hash = hexstrTo32Bytes(gjson.Get(data, "data.root").String())
	block.parentHash = hexstrTo32Bytes(gjson.Get(data, "data.header.message.parent_root").String())
	block.stateRoot = hexstrTo32Bytes(gjson.Get(data, "data.header.message.state_root").String())
	block.bodyRoot = hexstrTo32Bytes(gjson.Get(data, "data.header.message.body_root").String())

	return block
}

func hexstrTo32Bytes(hexString string) ([32]byte) {
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

func hexstrToUint64(hexString string) (uint64) {
	hex, _ := strconv.ParseInt(hexString, 0, 64)
	return uint64(hex)
}