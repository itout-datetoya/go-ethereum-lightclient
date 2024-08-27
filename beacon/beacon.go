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

	blockHeader.Slot = types.Slot(view.Uint64View(gjson.Get(data, "slot").Uint()))
	blockHeader.ProposerIndex = types.ValidatorIndex(view.Uint64View(gjson.Get(data, "proposer_index").Uint()))
	blockHeader.ParentRoot = tree.Root(util.HexstrTo32Bytes(gjson.Get(data, "parent_root").String()))
	blockHeader.StateRoot = tree.Root(util.HexstrTo32Bytes(gjson.Get(data, "state_root").String()))
	blockHeader.BodyRoot = tree.Root(util.HexstrTo32Bytes(gjson.Get(data, "body_root").String()))

	return blockHeader
}

func GetBeaconBlockHeader(slot uint64) (types.BeaconBlockHeader) {
	result := api.GetBeaconBlockHeader(slot)
	data := gjson.Get(result, "data.header.message").String()
	return ParseBeaconBlockHeader(data)
}