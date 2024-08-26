package types

import (
	"itout/go-ethereum-lightclient/configs"
	"github.com/protolambda/ztyp/tree"
	"github.com/protolambda/ztyp/view"
)

type BeaconBlockHeader struct {
	Slot          Slot           `json:"slot" yaml:"slot"`
	ProposerIndex ValidatorIndex `json:"proposer_index" yaml:"proposer_index"`
	ParentRoot    tree.Root           `json:"parent_root" yaml:"parent_root"`
	StateRoot     tree.Root           `json:"state_root" yaml:"state_root"`
	BodyRoot      tree.Root           `json:"body_root" yaml:"body_root"`
}

type Slot view.Uint64View
type Epoch view.Uint64View

type ValidatorIndex view.Uint64View

type SigningData struct {
	ObjectRoot tree.Root
	Domain Domain
}

type Domain []byte

type DomainType configs.BLSDomainType

type ForkData struct {
	CurrentVersion Version
	GenesisValidatorsRoot tree.Root
}

type Version []byte

func (b *BeaconBlockHeader) HashTreeRoot(hFn tree.HashFn) tree.Root {
	return hFn.HashTreeRoot(b.Slot, b.ProposerIndex, b.ParentRoot, b.StateRoot, b.BodyRoot)
}

func (s Slot) HashTreeRoot(hFn tree.HashFn) tree.Root {
	return view.Uint64View(s).HashTreeRoot(hFn)
}

func (i ValidatorIndex) HashTreeRoot(hFn tree.HashFn) tree.Root {
	return view.Uint64View(i).HashTreeRoot(hFn)
}

func (sData *SigningData) HashTreeRoot(hFn tree.HashFn) tree.Root {
	return hFn.HashTreeRoot(sData.ObjectRoot, sData.Domain)
}

func (domain Domain) HashTreeRoot(hFn tree.HashFn) tree.Root {
	return tree.Root(domain)
}

func (fData *ForkData) HashTreeRoot(hFn tree.HashFn) tree.Root {
	return hFn.HashTreeRoot(fData.CurrentVersion, fData.GenesisValidatorsRoot)
}

func (version Version) HashTreeRoot(hFn tree.HashFn) tree.Root {
	root := tree.Root{}
	copy(root[:], version[:])
	return root
}