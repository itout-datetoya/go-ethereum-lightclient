package helper

import (
	"itout/go-ethereum-lightclient/types"
	"github.com/protolambda/ztyp/tree"
)

func computeSigningRoot(b types.BeaconBlockHeader, domain types.Domain) tree.Root {
	hfn := tree.GetHashFn()
	sData := types.SigningData{ObjectRoot: b.HashTreeRoot(hfn), Domain: domain}

	return sData.HashTreeRoot(hfn)
}

func computeDomain(domainType types.DomainType, forkVersion types.Version, genesisValidatorsRoot tree.Root) types.Domain {
    forkDataRoot := computeForkDataRoot(forkVersion, genesisValidatorsRoot)
    return types.Domain(append(domainType, forkDataRoot[:28]...))
}

func computeForkDataRoot(currentVersion types.Version, genesisValidatorRoot tree.Root) tree.Root {
	hfn := tree.GetHashFn()
	fData := types.ForkData{CurrentVersion: currentVersion, GenesisValidatorsRoot: genesisValidatorRoot}
	return fData.HashTreeRoot(hfn)
}
