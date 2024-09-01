package beacon

import (
	"testing"
	//"fmt"

	"github.com/stretchr/testify/assert"
)

const EXE_URL_DEAULT = "https://mainnet.infura.io/v3/cdeb7402eca247e0a054717f350b4e50"
const BEACON_URL_DEFAULT = "https://docs-demo.quiknode.pro/eth/v1/beacon/"

const HASH = "0xacf662304a4dc3d37b112a5574c74343c8e6b30b0dfc7dcc16e4a494e53540ee"
const NUMBER = 20654229
const SLOT = 9862640
const BEACON_HASH = "0x6838b423f31b3a8a8147a340fdc0b16345ca60b89d5963f54e15169e6b50f503"

func TestBeacon(t *testing.T) {
	beaconBlockHeader := GetBeaconBlockHeader(SLOT, BEACON_URL_DEFAULT)

	assert.Equal(t, SLOT, int(beaconBlockHeader.Slot))
}