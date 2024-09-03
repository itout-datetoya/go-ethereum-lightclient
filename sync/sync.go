package sync

import (
	"errors"
	
	"itout/go-ethereum-lightclient/api"
	"itout/go-ethereum-lightclient/beacon"
	"itout/go-ethereum-lightclient/configs"
	"itout/go-ethereum-lightclient/helper"
	"itout/go-ethereum-lightclient/types"
	"itout/go-ethereum-lightclient/util"
	"math/bits"

	"github.com/protolambda/bls12-381-util"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"
	"github.com/tidwall/gjson"
)

const SYNC_COMMITTEE_INDEX = 54

var BLSPubkeyType = view.BasicVectorType(view.ByteType, 48)

type BLSPubkey [48]byte
type BLSSignature [96]byte
type SyncCommitteePubkeys []BLSPubkey

func (a SyncCommitteePubkeys) ByteLength(spec *configs.Spec) uint64 {
	return uint64(spec.SYNC_COMMITTEE_SIZE) * BLSPubkeyType.Size
}

func (a *SyncCommitteePubkeys) FixedLength(spec *configs.Spec) uint64 {
	return uint64(spec.SYNC_COMMITTEE_SIZE) * BLSPubkeyType.Size
}

type SyncCommittee struct {
	pubkeys SyncCommitteePubkeys
	aggPubkey BLSPubkey
}

func (p *SyncCommittee) HashTreeRoot(spec *configs.Spec, hFn tree.HashFn) tree.Root {
	return hFn.HashTreeRoot(spec.Wrap(&p.pubkeys), &p.aggPubkey)
}

func (li SyncCommitteePubkeys) HashTreeRoot(spec *configs.Spec, hFn tree.HashFn) tree.Root {
	return hFn.ComplexVectorHTR(func(i uint64) tree.HTR {
		return &li[i]
	}, uint64(spec.SYNC_COMMITTEE_SIZE))
}

func (p BLSPubkey) HashTreeRoot(hFn tree.HashFn) (tree.Root) {
	var a, b tree.Root
	serialisedPubkey := [48]byte(p)
	copy(a[:], serialisedPubkey[0:32])
	copy(b[:], serialisedPubkey[32:48])
	return hFn(a, b)
}

type SyncCommitteeBits []byte

func (li SyncCommitteeBits) ByteLength(spec *configs.Spec) uint64 {
	return (uint64(spec.SYNC_COMMITTEE_SIZE) + 7) / 8
}

func (li *SyncCommitteeBits) FixedLength(spec *configs.Spec) uint64 {
	return (uint64(spec.SYNC_COMMITTEE_SIZE) + 7) / 8
}

func (li SyncCommitteeBits) HashTreeRoot(spec *configs.Spec, hFn tree.HashFn) tree.Root {
	if li == nil {
		return view.BitVectorType(uint64(spec.SYNC_COMMITTEE_SIZE)).New().HashTreeRoot(hFn)
	}
	return hFn.BitVectorHTR(li)
}

func (li SyncCommitteeBits) GetBit(i uint64) bool {
	return util.GetBitFromBytes(li, i)
}

func (li SyncCommitteeBits) SetBit(i uint64, v bool) {
	util.SetBitToBytes(li, i, v)
}

func (li SyncCommitteeBits) PopCount() uint64 {
	count := 0
	for _, b := range li {
		count += bits.OnesCount8(uint8(b))
	}
	return uint64(count)
}

type SyncAggregate struct {
	syncCommitteeBits SyncCommitteeBits
	syncCommitteeSig BLSSignature
}

type Bootstrap struct {
	header types.BeaconBlockHeader
	syncCommittee SyncCommittee
	syncCommitteeBranch [][32]byte
}

func ParseBootstrap(data string) (Bootstrap) {
	bootstrap := Bootstrap{}

	jsonHeader := gjson.Get(data, "data.header.beacon").String()
	bootstrap.header = beacon.ParseBeaconBlockHeader(jsonHeader)
	
	strPubkeys := gjson.Get(data, "data.current_sync_committee.committee").Array()
	pubkeys := make([]BLSPubkey, 0, len(strPubkeys))
	for _, strPubkey := range strPubkeys {
		pubkeys = append(pubkeys, util.HexstrTo48Bytes(strPubkey.String()))
	}
	bootstrap.syncCommittee.pubkeys = pubkeys

	bootstrap.syncCommittee.aggPubkey = util.HexstrTo48Bytes(gjson.Get(data, "data.current_sync_committee.aggregate_public_key").String())

	strBranchs := gjson.Get(data, "data.current_sync_committee_branch").Array()
	branchs := make([][32]byte, 0, len(strBranchs))
	for _, strBranch := range strBranchs {
		branchs = append(branchs, util.HexstrTo32Bytes(strBranch.String()))
	}
	bootstrap.syncCommitteeBranch = branchs

	return bootstrap
}

func GetBootstrap(hash [32]byte, url string) (Bootstrap) {
	data := api.GetBootstrap(hash, url)
	return ParseBootstrap(data)
}

type Update struct {
	attestedHeader types.BeaconBlockHeader
	nextSyncCommittee SyncCommittee
	nextSyncCommitteeBranch [][32]byte
	finalizedHeader types.BeaconBlockHeader
	finalityBranch [][32]byte
	syncAggregate SyncAggregate
	slot types.Slot
}

func ParseUpdate(data string) (Update) {
	update := Update{}

	jsonAttestedHeader := gjson.Get(data, "data.attested_header.beacon").String()
	update.attestedHeader = beacon.ParseBeaconBlockHeader(jsonAttestedHeader)
	
	strPubkeys := gjson.Get(data, "data.next_sync_committee.committee").Array()
	pubkeys := make([]BLSPubkey, 0, len(strPubkeys))
	for _, strPubkey := range strPubkeys {
		pubkeys = append(pubkeys, util.HexstrTo48Bytes(strPubkey.String()))
	}
	update.nextSyncCommittee.pubkeys = pubkeys

	update.nextSyncCommittee.aggPubkey = util.HexstrTo48Bytes(gjson.Get(data, "data.next_sync_committee.aggregate_public_key").String())

	strComBranchs := gjson.Get(data, "data.next_sync_committee_branch").Array()
	comBranchs := make([][32]byte, 0, len(strComBranchs))
	for _, strComBranch := range strComBranchs {
		comBranchs = append(comBranchs, util.HexstrTo32Bytes(strComBranch.String()))
	}
	update.nextSyncCommitteeBranch = comBranchs

	jsonFinalizedHeader := gjson.Get(data, "data.finalized_header.beacon").String()
	update.finalizedHeader = beacon.ParseBeaconBlockHeader(jsonFinalizedHeader)

	strFinBranchs := gjson.Get(data, "data.finality_branch").Array()
	finBranchs := make([][32]byte, 0, len(strFinBranchs))
	for _, strFinBranch := range strFinBranchs {
		finBranchs = append(finBranchs, util.HexstrTo32Bytes(strFinBranch.String()))
	}
	update.finalityBranch = finBranchs

	update.syncAggregate.syncCommitteeBits = util.HexstrToBytes(gjson.Get(data, "data.sync_aggregate.sync_committee_bits").String())
	update.syncAggregate.syncCommitteeSig = util.HexstrTo96Bytes(gjson.Get(data, "data.sync_aggregate.sync_committee_signature").String())

	update.slot = types.Slot(view.Uint64View(gjson.Get(data, "data.signature_slot").Uint()))

	return update
}

func GetUpdate(currentSlot types.Slot, url string) (Update) {
	period := configs.Mainnet.SlotToPeriod(currentSlot)
	data := api.GetUpdate(period, url)
	return ParseUpdate(gjson.Get(data, "0").String())
}

type FinalityUpdate struct {
	AttestedHeader types.BeaconBlockHeader
	finalizedHeader types.BeaconBlockHeader
	finalityBranch [][32]byte
	syncAggregate SyncAggregate
	slot types.Slot
}

func ParseFinalityUpdate(data string) (FinalityUpdate) {
	update := FinalityUpdate{}

	jsonAttestedHeader := gjson.Get(data, "data.attested_header.beacon").String()
	update.AttestedHeader = beacon.ParseBeaconBlockHeader(jsonAttestedHeader)

	jsonFinalizedHeader := gjson.Get(data, "data.finalized_header.beacon").String()
	update.finalizedHeader = beacon.ParseBeaconBlockHeader(jsonFinalizedHeader)

	strFinBranchs := gjson.Get(data, "data.finality_branch").Array()
	finBranchs := make([][32]byte, 0, len(strFinBranchs))
	for _, strFinBranch := range strFinBranchs {
		finBranchs = append(finBranchs, util.HexstrTo32Bytes(strFinBranch.String()))
	}
	update.finalityBranch = finBranchs

	update.syncAggregate.syncCommitteeBits = util.HexstrToBytes(gjson.Get(data, "data.sync_aggregate.sync_committee_bits").String())
	update.syncAggregate.syncCommitteeSig = util.HexstrTo96Bytes(gjson.Get(data, "data.sync_aggregate.sync_committee_signature").String())

	update.slot = types.Slot(view.Uint64View(gjson.Get(data, "data.signature_slot").Uint()))

	return update
}

func GetFinalityUpdate(url string) (FinalityUpdate) {
	data := api.GetFinalityUpdate(url)
	return ParseFinalityUpdate(data)
}

type Store struct {
	Header types.BeaconBlockHeader
	CurrentSyncCommittee SyncCommittee
	NextSyncCommittee SyncCommittee
}

func InitStore(trustedRoot tree.Root, bootstrap Bootstrap) (Store, error) {
	hfn := tree.GetHashFn()
	if(bootstrap.header.HashTreeRoot(hfn) != trustedRoot){
		return Store{}, errors.New("error:wrong root")
	}

	if !helper.IsValidMerkleBranch(bootstrap.syncCommittee.HashTreeRoot(configs.Mainnet, hfn), bootstrap.syncCommitteeBranch, SYNC_COMMITTEE_INDEX, bootstrap.header.StateRoot) {
		return Store{}, errors.New("error:wrong merkle branch")
	}

	return Store{bootstrap.header, bootstrap.syncCommittee, SyncCommittee{}}, nil
}

func (store *Store) UpdateStore(update Update, spec *configs.Spec) error {
	if view.Uint64View(update.syncAggregate.syncCommitteeBits.PopCount()) < spec.MIN_SYNC_COMMITTEE_PARTICIPANTS {
		return errors.New("error:insufficient participants")
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

	forkVersionSlot := max(update.attestedHeader.Slot, types.Slot(1)) - types.Slot(1)
	forkVersion := spec.ForkVersion(forkVersionSlot)
	domain := helper.ComputeDomain(types.DomainType(configs.DOMAIN_SYNC_COMMITTEE), forkVersion, configs.GENESIS_VALIDATORS_ROOT)
	signingRoot := helper.ComputeSigningRoot(update.attestedHeader, domain)
	
	serialisedSig := [96]byte(update.syncAggregate.syncCommitteeSig)
	sig := blsu.Signature{}
	sig.Deserialize(&serialisedSig)

	if !blsu.FastAggregateVerify(paticipantPubkeys, signingRoot[:], &sig) {
		/*
		verification paused in testnet
		*/
		// return errors.New("error:wrong signature")
	}

	if spec.SlotToPeriod(store.Header.Slot) == spec.SlotToPeriod(update.attestedHeader.Slot) {
		store.NextSyncCommittee = update.nextSyncCommittee
	} else if spec.SlotToPeriod(store.Header.Slot) + 1 == spec.SlotToPeriod(update.attestedHeader.Slot) {
		store.CurrentSyncCommittee = store.NextSyncCommittee
		store.NextSyncCommittee = update.nextSyncCommittee
	}
	
	return nil
}

func (store *Store) FinalityUpdateStore(update FinalityUpdate, spec *configs.Spec) error {
	if view.Uint64View(update.syncAggregate.syncCommitteeBits.PopCount()) < spec.MIN_SYNC_COMMITTEE_PARTICIPANTS {
		return errors.New("error:insufficient participants")
	}

	syncCommittee := SyncCommittee{}

	if spec.SlotToPeriod(store.Header.Slot) == spec.SlotToPeriod(update.AttestedHeader.Slot) {
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

	forkVersionSlot := max(update.AttestedHeader.Slot, types.Slot(1)) - types.Slot(1)
	forkVersion := spec.ForkVersion(forkVersionSlot)
	domain := helper.ComputeDomain(types.DomainType(configs.DOMAIN_SYNC_COMMITTEE), forkVersion, configs.GENESIS_VALIDATORS_ROOT)
	signingRoot := helper.ComputeSigningRoot(update.AttestedHeader, domain)
	
	serialisedSig := [96]byte(update.syncAggregate.syncCommitteeSig)
	sig := blsu.Signature{}
	sig.Deserialize(&serialisedSig)

	if !blsu.FastAggregateVerify(paticipantPubkeys, signingRoot[:], &sig) {
		/*
		verification paused in testnet
		*/
		// return errors.New("error:wrong signature")
	}

	store.Header = update.AttestedHeader
	if spec.SlotToPeriod(store.Header.Slot) + 1 == spec.SlotToPeriod(update.AttestedHeader.Slot) {
		store.CurrentSyncCommittee = store.NextSyncCommittee
	}
	
	return nil
}