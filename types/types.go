package types

type SigningData struct {
	objectRoot [32]byte
	domain Domain
}

type Domain [32]byte

type DomainType [4]byte

type ForkData struct {
	currentVersion [4]byte
	genesisValidatorsRoot [32]byte
}