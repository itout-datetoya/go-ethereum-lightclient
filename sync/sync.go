package sync

import (
	"fmt"
	"errors"
	"itout/go-ethereum-lightclient/util"
	"itout/go-ethereum-lightclient/types"
	"itout/go-ethereum-lightclient/rpc"
	"itout/go-ethereum-lightclient/helper"
	"itout/go-ethereum-lightclient/configs"
	"github.com/tidwall/gjson"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"
	"github.com/protolambda/ztyp/codec"
	"github.com/protolambda/bls12-381-util"
)

const SYNC_COMMITTEE_INDEX = 54

var BLSPubkeyType = view.BasicVectorType(view.ByteType, 48)

type BLSPubkey [48]byte
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

type SyncAggregate struct {
	syncCommitteeBits []bool
	syncCommitteeSig blsu.Signature
}

type Bootstrap struct {
	header types.BeaconBlockHeader
	syncCommittee SyncCommittee
	syncCommitteeBranch [][32]byte
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

type Store struct {
	header types.BeaconBlockHeader
	currentSyncCommittee SyncCommittee
	nextSyncCommittee SyncCommittee
}

func InitStore(trustedRoot tree.Root, bootstrap Bootstrap) (Store, error) {
	hfn := tree.GetHashFn()
	if(bootstrap.header.HashTreeRoot(hfn) != trustedRoot){
		return Store{}, errors.New("Error:wrong root")
	}

	if !helper.IsValidMerkleBranch(bootstrap.syncCommittee.HashTreeRoot(configs.Mainnet, hfn), bootstrap.syncCommitteeBranch, SYNC_COMMITTEE_INDEX, bootstrap.header.StateRoot) {
		return Store{}, errors.New("Error:wrong merkle branch")
	}

	return Store{bootstrap.header, bootstrap.syncCommittee, SyncCommittee{}}, nil
}