package helper

import (
	"crypto/sha256"
	"itout/go-ethereum-lightclient/types"
	"github.com/protolambda/ztyp/tree"
)

func ComputeSigningRoot(b types.BeaconBlockHeader, domain types.Domain) tree.Root {
	hfn := tree.GetHashFn()
	sData := types.SigningData{ObjectRoot: b.HashTreeRoot(hfn), Domain: domain}

	return sData.HashTreeRoot(hfn)
}

func ComputeDomain(domainType types.DomainType, forkVersion types.Version, genesisValidatorsRoot tree.Root) types.Domain {
    forkDataRoot := ComputeForkDataRoot(forkVersion, genesisValidatorsRoot)
    return types.Domain(append(domainType[:], forkDataRoot[:28]...))
}

func ComputeForkDataRoot(currentVersion types.Version, genesisValidatorRoot tree.Root) tree.Root {
	hfn := tree.GetHashFn()
	fData := types.ForkData{CurrentVersion: currentVersion, GenesisValidatorsRoot: genesisValidatorRoot}
	return fData.HashTreeRoot(hfn)
}

func IsValidMerkleBranch(leaf [32]byte, branch [][32]byte, index uint64, root tree.Root) (bool) {
	hasher := sha256.New()
	for _, sibling := range branch {
		hasher.Reset()
		if index&1 == 0 {
			hasher.Write(leaf[:])
			hasher.Write(sibling[:])
		} else {
			hasher.Write(sibling[:])
			hasher.Write(leaf[:])
		}
		hasher.Sum(leaf[:0])
		if index >>= 1; index == 0 {
			return false
		}
	}
	if index != 1 {
		return false
	}
	if tree.Root(leaf) != root {
		return false
	}
	return true
}
