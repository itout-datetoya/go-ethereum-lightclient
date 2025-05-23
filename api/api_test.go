package api

import (
	"itout/go-ethereum-lightclient/util"
	"itout/go-ethereum-lightclient/configs"
	"itout/go-ethereum-lightclient/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const EXE_URL_DEAULT = "https://mainnet.infura.io/v3/cdeb7402eca247e0a054717f350b4e50"
const BEACON_URL_DEFAULT = "https://eth-holesky-beacon.blastapi.io/6cc7d5ce-13ed-4621-84b4-da82d9a464df/eth/v1/beacon/"

const HASH = "0xacf662304a4dc3d37b112a5574c74343c8e6b30b0dfc7dcc16e4a494e53540ee"
const NUMBER = 20654229
const SLOT = 2440637
const BEACON_HASH = "0x1377b7ee2517553b7ebe444cb4677fb76aa9b29b407b38799b8a5d5e293cd5e0"

func TestGetBlock(t *testing.T) {
	hash := util.HexstrTo32Bytes(HASH)
	dataByHash := GetBlockByHash(hash, EXE_URL_DEAULT)
	dataByNumber := GetBlockByNumber(NUMBER, EXE_URL_DEAULT)

	assert.Equal(t, HASH, gjson.Get(dataByHash, "result.hash").String())
	assert.Equal(t, HASH, gjson.Get(dataByNumber, "result.hash").String())
}

func TestGetBeaconBlockHeader(t *testing.T) {
	data := GetBeaconBlockHeader(SLOT, BEACON_URL_DEFAULT)

	assert.Equal(t, BEACON_HASH, gjson.Get(data, "data.0.root").String())
}

func TestGetBootstrap(t *testing.T) {
	hash := util.HexstrTo32Bytes(BEACON_HASH)
	data := GetBootstrap(hash, BEACON_URL_DEFAULT)

	assert.Equal(t, configs.Mainnet.SlotToPeriod(SLOT), configs.Mainnet.SlotToPeriod(types.Slot(gjson.Get(data, "data.header.beacon.slot").Int())))
}

func TestGetUpdate(t *testing.T) {
	period := configs.Mainnet.SlotToPeriod(SLOT)
	data := GetUpdate(period, BEACON_URL_DEFAULT)

	assert.Equal(t, period, configs.Mainnet.SlotToPeriod(types.Slot(gjson.Get(data, "0.data.attested_header.beacon.slot").Int())))
}

func TestGetFinalityUpdate(t *testing.T) {
	data := GetFinalityUpdate(BEACON_URL_DEFAULT)

	assert.NotZero(t, gjson.Get(data, "data.attested_header.beacon.slot").Int())
}