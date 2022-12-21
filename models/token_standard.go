package models

// TokenStandard represent the interface or standard that
// a token is respecting. Each chain may have their own token standards.
type TokenStandard string

const (
	// -- Ethereum/EVM standards

	// TokenStandardERC721 is https://eips.ethereum.org/EIPS/eip-721
	TokenStandardERC721 TokenStandard = "erc721"
	// TokenStandardERC1155 is https://eips.ethereum.org/EIPS/eip-1155
	TokenStandardERC1155 TokenStandard = "erc1155"

	// -- Solana token standards

	// TokenStandardMetaplexNonFungible is https://docs.metaplex.com/programs/token-metadata/token-standard
	TokenStandardMetaplexNonFungible TokenStandard = "metaplex_non_fungible"
)
