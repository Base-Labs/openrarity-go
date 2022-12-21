package scoring

import (
	"github.com/Base-Labs/openrarity/models"
	"github.com/pkg/errors"
)

// IScorer is the main class to score rarity scores for a given
// collection and token(s) based on the default OpenRarity scoring
// algorithm.
type IScorer interface {
	IScoreHandler
	// ValidateCollection is used to validate collection eligibility for OpenRarity scoring
	ValidateCollection(collection models.ICollection) error
	// ScoreCollection is used to score all tokens on collection.tokens
	ScoreCollection(collection models.ICollection) ([]float64, error)
	// ScoreCollections is used to score all tokens in every collection provided.
	ScoreCollections(collection []models.ICollection) ([][]float64, error)
}

// IScoreHandler class is an interface for different scoring algorithms to
// implement. Subclasses are responsible to ensure the batch functions are
// efficient for their particular algorithm.
type IScoreHandler interface {
	// ScoreToken is used to score an individual token based on the traits' distribution across
	// the whole collection.
	ScoreToken(collection models.ICollection, token models.IToken) (float64, error)
	// ScoreTokens should be used if you only want to score a batch of tokens that belong to collection.
	// This will typically be more efficient than calling score_token for each
	// token in `tokens`.
	ScoreTokens(collection models.ICollection, tokens []models.IToken) ([]float64, error)
}

// Scorer is the main class to score rarity scores for a given
// collection and token(s) based on the default OpenRarity scoring
// algorithm.
type Scorer struct {
	handler IScoreHandler
}

var _ IScorer = &Scorer{}

// NewScorer is the constructor of Scorer
func NewScorer(handler IScoreHandler) *Scorer {
	return &Scorer{
		handler: handler,
	}
}

// ValidateCollection is used to validate collection eligibility for OpenRarity scoring
func (c *Scorer) ValidateCollection(collection models.ICollection) error {
	if collection.HasNumericAttribute() {
		return errors.New("OpenRarity currently does not support collections with " +
			"numeric or date traits")
	}
	allowedStandards := []models.TokenStandard{
		models.TokenStandardERC721,
		models.TokenStandardMetaplexNonFungible,
	}
	if !models.IsSubset(allowedStandards, collection.TokenStandards()) {
		return errors.New("OpenRarity currently only supports ERC721/Non-fungible standards")
	}
	return nil
}

// ScoreToken is used to score an individual token based on the
// traits' distribution across the whole collection.
func (c *Scorer) ScoreToken(collection models.ICollection, token models.IToken) (float64, error) {
	if err := c.ValidateCollection(collection); err != nil {
		return 0, err
	}
	return c.handler.ScoreToken(collection, token)
}

// ScoreTokens should be used if you only want to score a batch of tokens that belong to collection.
// This will typically be more efficient than calling score_token for each
// token in `tokens`.
func (c *Scorer) ScoreTokens(collection models.ICollection, tokens []models.IToken) ([]float64, error) {
	if err := c.ValidateCollection(collection); err != nil {
		return nil, err
	}
	return c.handler.ScoreTokens(collection, tokens)
}

// ScoreCollection is used to score all tokens on collection.tokens
func (c *Scorer) ScoreCollection(collection models.ICollection) ([]float64, error) {
	if err := c.ValidateCollection(collection); err != nil {
		return nil, err
	}
	return c.handler.ScoreTokens(collection, collection.Tokens())
}

// ScoreCollections is used to score all tokens in every collection provided.
func (c *Scorer) ScoreCollections(collections []models.ICollection) ([][]float64, error) {
	allScores := make([][]float64, 0, len(collections))
	for _, collection := range collections {
		if err := c.ValidateCollection(collection); err != nil {
			return nil, err
		}
		scores, err := c.ScoreTokens(collection, collection.Tokens())
		if err != nil {
			return nil, err
		}
		allScores = append(allScores, scores)
	}
	return allScores, nil
}
