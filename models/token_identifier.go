package models

// IdentifierType defines the type of identifier
type IdentifierType string

// defines a set of identifier types
const (
	IdentifierTypeEVMContract       IdentifierType = "evm_contract"
	IdentifierTypeSolanaMintAddress IdentifierType = "solana_mint_address"
)

// ITokenIdentifier is used to specify how the collection is identified and the
// logic used to group the NFTs together
type ITokenIdentifier interface {
	// IdentifierType is used to obtain the identifier type of the current Token.
	IdentifierType() IdentifierType
}

// EVMContractTokenIdentifier indicates that this token is identified by the contract address and token ID number.
// This identifier is based off of the interface as defined by ERC721 and ERC1155,
// where unique tokens belong to the same contract but have their own numeral token id.
type EVMContractTokenIdentifier struct {
	contractAddress string
	tokenID         int
}

var _ ITokenIdentifier = &EVMContractTokenIdentifier{}

// NewEVMContractTokenIdentifier is the constructor of EVMContractTokenIdentifier
func NewEVMContractTokenIdentifier(
	contractAddress string,
	tokenID int,
) EVMContractTokenIdentifier {
	return EVMContractTokenIdentifier{
		contractAddress: contractAddress,
		tokenID:         tokenID,
	}
}

// IdentifierType is used to obtain the identifier type of the current Token.
func (c EVMContractTokenIdentifier) IdentifierType() IdentifierType {
	return IdentifierTypeEVMContract
}

// SolanaMintAddressTokenIdentifier indicates that this token is identified by their solana account address.
// This identifier is based off of the interface defined by the Solana SPL token
// standard where every such token is declared by creating a mint account.
type SolanaMintAddressTokenIdentifier struct {
	mintAddress string
}

var _ ITokenIdentifier = &SolanaMintAddressTokenIdentifier{}

// NewSolanaMintAddressTokenIdentifier is the constructor of SolanaMintAddressTokenIdentifier
func NewSolanaMintAddressTokenIdentifier(mintAddress string) SolanaMintAddressTokenIdentifier {
	return SolanaMintAddressTokenIdentifier{
		mintAddress: mintAddress,
	}
}

// IdentifierType is used to obtain the identifier type of the current Token.
func (c SolanaMintAddressTokenIdentifier) IdentifierType() IdentifierType {
	return IdentifierTypeSolanaMintAddress
}
