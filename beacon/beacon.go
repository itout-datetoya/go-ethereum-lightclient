package beacon

import (
	"itout/go-ethereum-lightclient/util"
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

	block.slot = util.HexstrToUint64(gjson.Get(data, "data.header.message.slot").String())
	block.hash = util.HexstrTo32Bytes(gjson.Get(data, "data.root").String())
	block.parentHash = util.HexstrTo32Bytes(gjson.Get(data, "data.header.message.parent_root").String())
	block.stateRoot = util.HexstrTo32Bytes(gjson.Get(data, "data.header.message.state_root").String())
	block.bodyRoot = util.HexstrTo32Bytes(gjson.Get(data, "data.header.message.body_root").String())

	return block
}