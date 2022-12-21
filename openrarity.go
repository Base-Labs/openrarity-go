package openrarity

import (
	"github.com/Base-Labs/openrarity/models"
	"github.com/Base-Labs/openrarity/scoring"
	"github.com/Base-Labs/openrarity/scoring/handlers"
)

// NewOpenRarityScorer is used to build the default scorer of open rarity.
func NewOpenRarityScorer() scoring.IScorer {
	return scoring.NewScorer(handlers.NewInformationContentScoringHandler())
}

// export a set of types
type (
	IToken                = models.IToken
	ICollection           = models.ICollection
	ITokenRarity          = models.ITokenRarity
	IStringAttribute      = models.IStringAttribute
	IAttribute            = models.IAttribute
	ITokenMetadata        = models.ITokenMetadata
	ITokenRankingFeatures = models.ITokenRankingFeatures
	ITokenIdentifier      = models.ITokenIdentifier
)

// export a set of methods
var (
	// NewCollection is the constructor of Collection
	NewCollection = models.NewCollection
	// NewERC721Token Creates a Token class representing an ERC721 evm token given the following
	// parameters.
	NewERC721Token = models.NewERC721Token
	// NewToken is the constructor of Token
	NewToken = models.NewToken
)
