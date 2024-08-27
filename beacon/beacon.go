package beacon

import (
	"itout/go-ethereum-lightclient/util"
	"itout/go-ethereum-lightclient/types"
	"itout/go-ethereum-lightclient/api"
	"github.com/tidwall/gjson"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"
)

type BeaconBlock struct {
	slot uint64
	root [32]byte
	blockheader types.BeaconBlockHeader
}

func ParseBeaconBlockHeader(data string) (types.BeaconBlockHeader) {
	blockHeader := types.BeaconBlockHeader{}

	blockHeader.Slot = types.Slot(view.Uint64View(util.HexstrToUint64(gjson.Get(data, "data.header.message.slot").String())))
	blockHeader.ProposerIndex = types.ValidatorIndex(view.Uint64View(util.HexstrToUint64(gjson.Get(data, "data.header.message.proposer_index").String())))
	blockHeader.ParentRoot = tree.Root(util.HexstrTo32Bytes(gjson.Get(data, "data.header.message.parent_root").String()))
	blockHeader.StateRoot = tree.Root(util.HexstrTo32Bytes(gjson.Get(data, "data.header.message.state_root").String()))
	blockHeader.BodyRoot = tree.Root(util.HexstrTo32Bytes(gjson.Get(data, "data.header.message.body_root").String()))

	return blockHeader
}

func GetBeaconBlockHeader(slot uint64) (types.BeaconBlockHeader) {
	data := api.GetBeaconBlockHeader(slot)
	return ParseBeaconBlockHeader(data)
}