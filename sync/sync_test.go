package sync

import (
	"encoding/hex"
	"itout/go-ethereum-lightclient/configs"
	"itout/go-ethereum-lightclient/types"
	"itout/go-ethereum-lightclient/util"
	"itout/go-ethereum-lightclient/helper"
	"testing"
	"errors"
	"fmt"
	"io"
	"os"
	"crypto/rand"
	"github.com/protolambda/ztyp/view"
	"github.com/protolambda/bls12-381-util"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const BEACON_URL_DEFAULT = "https://eth-holesky-beacon.blastapi.io/6cc7d5ce-13ed-4621-84b4-da82d9a464df/eth/v1/beacon/"

const SLOT = 2452759
const BEACON_HASH = "0x4c1282286c0c6630886dc6961e3beef0dd50f886bd029ae09a07a7c7aa89ec02"

var BLANCH = []string{
	"0x88b104b329048db035143a2905ea865a1b2b9e5981db206f0f030a00a4e5c8e0",
	"0x45a7431e713880455807bbb305a73a657f5f827b1db35241ab4429743b8d998b",
	"0xa74e8782b41bd0cf7ba7aec10640e5994eda1282a78d302ea90abf17c0a444df",
	"0x3901f689c98dc6745d54d4922605c7828b0f297483c0715661ba05148345c6a8",
	"0xb63b1856f967249209910bcd9d5323dee988cc3d89610f0296f6d67cf0d43247",
}

var FIN_BLANCH = []string{
	"0xfe2a010000000000000000000000000000000000000000000000000000000000",
	"0xfa4c9d948cbf6bea02453894e4506c11ad8f515a71a093d1b604de3b8ad817a1",
	"0x85057d0f56902edfb6c3b374585645e7ec36dd384c5479ae26b234d3a712b8d7",
	"0xbfafa032d91756368c199560f2786d0b09fd90c409c3d72d087e6dfc602738f6",
	"0x67645a3b02b55b49d01c623d0d8ec6ef340ef32f13ffee075ba739212acbb158",
	"0x962f86effc1f2404dccd0101ce2e8772350511170cbdfa1e5a111db25efcddc4",
}

const HOLESKY_GENESIS_VALIDATORS_ROOT = "0x9143aa7c615a7f7115e2b6aac319c03529df8242ae705fba9df39b79c59fa8b1"


func TestGetBootstrap(t *testing.T) {
	bootstrap := GetBootstrap(util.HexstrTo32Bytes(BEACON_HASH), BEACON_URL_DEFAULT)

	assert.Equal(t, len(BLANCH), len(bootstrap.syncCommitteeBranch))

	for i, branch := range bootstrap.syncCommitteeBranch {
		assert.Equal(t, BLANCH[i], "0x" + hex.EncodeToString(branch[:]))
	}
}

func TestGetUpdate(t *testing.T) {
	update := GetUpdate(types.Slot(SLOT), BEACON_URL_DEFAULT)

	assert.Equal(t, len(FIN_BLANCH), len(update.finalityBranch))
	for i, branch := range update.finalityBranch {
		fmt.Println("0x" + hex.EncodeToString(branch[:]))
		assert.Equal(t, FIN_BLANCH[i], "0x" + hex.EncodeToString(branch[:]))
	}
}

func TestGetFinalityUpdate(t *testing.T) {
	update := GetFinalityUpdate(BEACON_URL_DEFAULT)

	assert.NotZero(t, len(update.finalityBranch))
}

func TestInitStore(t *testing.T) {
	bootstrap := GetBootstrap(util.HexstrTo32Bytes(BEACON_HASH), BEACON_URL_DEFAULT)

	store, err := InitStore(util.HexstrTo32Bytes(BEACON_HASH), bootstrap)
	if err != nil {
		fmt.Println("Client failed to start")
	}
	assert.Equal(t, int(bootstrap.header.Slot), int(store.Header.Slot))
}

func TestAggSig(t *testing.T) {
	randomBytes := make([]byte, 32)
	pubkeys := make([]*blsu.Pubkey, 0, 512)
	paticipantPubkeys := make([]*blsu.Pubkey, 0, 512)
	sigs := make([]*blsu.Signature, 0, 512)
	paticipantSigs := make([]*blsu.Signature, 0, 512)
	signingRoot := util.HexstrTo32Bytes(BEACON_HASH)

	for i := 0; i < 512; i++ {
		rand.Read(randomBytes)
		sk := blsu.SecretKey{}
		sk.Deserialize((*[32]byte)(randomBytes))
		pk, _ := blsu.SkToPk(&sk)
		pubkeys = append(pubkeys, pk)
		sig := blsu.Sign(&sk, signingRoot[:])
		sigs = append(sigs, sig)
	}

	syncCommitteeBits := SyncCommitteeBits(util.HexstrToBytes("0x4bf30e1cd4c95876eba6ae4eafc1ef3f997a66596b56862beede6477abf15d1278bd729f3b76bb9f46871109437c75fecfa37dbc6ad36db78b4d5e9bef877c3c"))
	for i, pubkey := range pubkeys {
		if syncCommitteeBits.GetBit(uint64(i)) {
			paticipantPubkeys = append(paticipantPubkeys, pubkey)
			paticipantSigs = append(paticipantSigs, sigs[i])
		}
	}
	tempSig := paticipantSigs[0]
	paticipantSigs[0] = paticipantSigs[1]
	paticipantSigs[1] = tempSig

	aggSig, _ := blsu.Aggregate(paticipantSigs)

	if !blsu.FastAggregateVerify(paticipantPubkeys, signingRoot[:], aggSig) {
		fmt.Println("error:wrong signature")
	} else {
		fmt.Println("valid signature")
	}
}

func TestUpdateStore(t *testing.T) {

	json298, _ := os.Open("update_298.json")
	defer json298.Close()
	jsonData298, _ := io.ReadAll(json298)
	update_298 := ParseUpdate(gjson.Get(string(jsonData298), "0").String())

	store := Store{Header: update_298.attestedHeader, CurrentSyncCommittee: update_298.nextSyncCommittee, NextSyncCommittee:  update_298.nextSyncCommittee}

	json299, _ := os.Open("update_299.json")
	defer json299.Close()
	jsonData299, _ := io.ReadAll(json299)
	update_299 := ParseUpdate(gjson.Get(string(jsonData299), "0").String())
	
	err := store.UpdateStoreHolesky(update_299, configs.Mainnet)
	if err != nil {
		fmt.Println(err)
	}
	
	assert.Equal(t, int(update_299.attestedHeader.Slot), int(store.Header.Slot))
}

func (store *Store) UpdateStoreHolesky(update Update, spec *configs.Spec) error {
	if view.Uint64View(update.syncAggregate.syncCommitteeBits.PopCount()) < spec.MIN_SYNC_COMMITTEE_PARTICIPANTS {
		return errors.New("error:insufficient participants")
	}

	if store.Header.Slot >= update.attestedHeader.Slot {
		return errors.New("error:previous attested header")
	}

	syncCommittee := SyncCommittee{}

	if spec.SlotToPeriod(store.Header.Slot) == spec.SlotToPeriod(update.attestedHeader.Slot) {
		syncCommittee = store.CurrentSyncCommittee
	} else {
		syncCommittee = store.NextSyncCommittee
	}

	paticipantPubkeys := make([]*blsu.Pubkey, 0, len(syncCommittee.pubkeys))
	for i, data := range syncCommittee.pubkeys {
		if update.syncAggregate.syncCommitteeBits.GetBit(uint64(i)) {
			serialisedPubkey := [48]byte(data)
			pubkey := blsu.Pubkey{}
			pubkey.Deserialize(&serialisedPubkey)
			paticipantPubkeys = append(paticipantPubkeys, &pubkey)
		}
	}

	
	forkVersion := spec.DENEB_FORK_VERSION
	domain := helper.ComputeDomain(types.DomainType(configs.DOMAIN_SYNC_COMMITTEE), forkVersion, util.HexstrTo32Bytes(HOLESKY_GENESIS_VALIDATORS_ROOT))
	signingRoot := helper.ComputeSigningRoot(update.attestedHeader, domain)
	
	serialisedSig := [96]byte(update.syncAggregate.syncCommitteeSig)
	sig := blsu.Signature{}
	sig.Deserialize(&serialisedSig)

	if !blsu.FastAggregateVerify(paticipantPubkeys, signingRoot[:], &sig) {
		fmt.Println("error:wrong signature")
	}

	store.Header = update.attestedHeader
	if spec.SlotToPeriod(store.Header.Slot) == spec.SlotToPeriod(update.attestedHeader.Slot) {
		store.NextSyncCommittee = update.nextSyncCommittee
	} else if spec.SlotToPeriod(store.Header.Slot) + 1 == spec.SlotToPeriod(update.attestedHeader.Slot) {
		store.CurrentSyncCommittee = store.NextSyncCommittee
		store.NextSyncCommittee = update.nextSyncCommittee
	}
	
	return nil
}

func TestFinalityUpdateStore(t *testing.T) {
	bootstrap := GetBootstrap(util.HexstrTo32Bytes(BEACON_HASH), BEACON_URL_DEFAULT)

	store, _ := InitStore(util.HexstrTo32Bytes(BEACON_HASH), bootstrap)

	update := GetFinalityUpdate(BEACON_URL_DEFAULT)

	fmt.Println("Slot: ", int(store.Header.Slot))
	err := store.FinalityUpdateStore(update, configs.Mainnet)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Slot: ", int(store.Header.Slot))

	assert.Equal(t, int(update.AttestedHeader.Slot), int(store.Header.Slot))
}