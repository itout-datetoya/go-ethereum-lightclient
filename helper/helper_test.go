package helper

import (
	"testing"
	"fmt"
	"encoding/hex"
	"itout/go-ethereum-lightclient/beacon"
	"itout/go-ethereum-lightclient/types"
	"itout/go-ethereum-lightclient/configs"
	"github.com/protolambda/ztyp/tree"
	"github.com/stretchr/testify/assert"
)

const BEACON_URL_DEFAULT = "https://docs-demo.quiknode.pro/eth/v1/beacon/"

const SLOT = 9872124
const BEACON_HASH = "0x8eb5c9952474e9fe270e5a1823a198d2b8a5112023306dbef34418484ca7d458"

func TestComputeSigningRoot(t *testing.T) {
	attestedHeader := beacon.GetBeaconBlockHeader(SLOT, BEACON_URL_DEFAULT)

	forkVersionSlot := max(types.Slot(SLOT), types.Slot(1)) - types.Slot(1)
	forkVersion := configs.Mainnet.ForkVersion(forkVersionSlot)
	domain := ComputeDomain(types.DomainType(configs.DOMAIN_SYNC_COMMITTEE), forkVersion, configs.GENESIS_VALIDATORS_ROOT)
	ComputeSigningRoot(attestedHeader, domain)

	root := attestedHeader.HashTreeRoot(tree.GetHashFn())
	strRoot := hex.EncodeToString(root[:])
	fmt.Println("Root: ", strRoot)

	assert.Equal(t, BEACON_HASH, "0x" + strRoot)
}