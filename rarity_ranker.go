package openrarity

import (
	"math"
	"sort"

	"github.com/Base-Labs/openrarity/models"
	"github.com/Base-Labs/openrarity/scoring"
	"github.com/pkg/errors"
)

// IRarityRanker is used to rank a set of tokens given their rarity scores.
type IRarityRanker interface {
	// RankCollection is used to rank tokens in the collection with the default scorer implementation.
	// Scores that are higher indicate a higher rarity, and thus a lower rank.
	//
	// Tokens with the same score will be assigned the same rank, e.g. we use RANK
	// (vs. DENSE_RANK).
	// Example: 1, 2, 2, 2, 5.
	// Scores are considered the same rank if they are within about 9 decimal digits
	// of each other.
	RankCollection(collection models.ICollection, scorer scoring.IScorer) ([]models.ITokenRarity, error)
	// SetRarityRanks is used to rank a set of tokens according to OpenRarity algorithm.
	// To account for additional factors like unique items in a collection,
	// OpenRarity implements multifactorial sort. Current sort algorithm uses two
	// factors: unique attributes count and Information Content score, in order.
	// Tokens with the same score will be assigned the same rank, e.g. we use RANK
	// (vs. DENSE_RANK).
	// Example: 1, 2, 2, 2, 5.
	// Scores are considered the same rank if they are within about 9 decimal digits
	// of each other.
	SetRarityRanks(tokenRarities []models.ITokenRarity) ([]models.ITokenRarity, error)
}

// RarityRanker is used to rank a set of tokens given their rarity scores.
type RarityRanker struct{}

var _ IRarityRanker = &RarityRanker{}

// NewRarityRanker is the constructor of RarityRanker
func NewRarityRanker() *RarityRanker {
	return &RarityRanker{}
}

// RankCollection is used to rank tokens in the collection with the default scorer implementation.
// Scores that are higher indicate a higher rarity, and thus a lower rank.
//
// Tokens with the same score will be assigned the same rank, e.g. we use RANK
// (vs. DENSE_RANK).
// Example: 1, 2, 2, 2, 5.
// Scores are considered the same rank if they are within about 9 decimal digits
// of each other.
func (c *RarityRanker) RankCollection(collection models.ICollection, scorer scoring.IScorer) ([]models.ITokenRarity, error) {
	if collection == nil || collection.Tokens() == nil {
		return nil, nil
	}
	tokens := collection.Tokens()
	scores, err := scorer.ScoreTokens(collection, tokens)
	if err != nil {
		return nil, err
	}
	if len(tokens) != len(scores) {
		return nil, errors.New("dimension of scores doesn't match dimension of tokens")
	}
	tokenRarities := make([]models.ITokenRarity, 0, len(tokens))
	for idx, token := range tokens {
		tokenFeatures := models.ExtractUniqueAttributeCount(token, collection)
		tokenRarities = append(tokenRarities,
			models.NewTokenRarity(token, scores[idx], tokenFeatures),
		)
	}
	return c.SetRarityRanks(tokenRarities)
}

// SetRarityRanks is used to rank a set of tokens according to OpenRarity algorithm.
// To account for additional factors like unique items in a collection,
// OpenRarity implements multifactorial sort. Current sort algorithm uses two
// factors: unique attributes count and Information Content score, in order.
// Tokens with the same score will be assigned the same rank, e.g. we use RANK
// (vs. DENSE_RANK).
// Example: 1, 2, 2, 2, 5.
// Scores are considered the same rank if they are within about 9 decimal digits
// of each other.
func (c *RarityRanker) SetRarityRanks(tokenRarities []models.ITokenRarity) ([]models.ITokenRarity, error) {
	sort.Slice(tokenRarities, func(i, j int) bool {
		if delta := tokenRarities[i].TokenFeatures().UniqueAttributeCount() -
			tokenRarities[j].TokenFeatures().UniqueAttributeCount(); delta != 0 {
			return delta > 0
		}
		return tokenRarities[i].Score() > tokenRarities[j].Score()
	})
	for i, tokenRarity := range tokenRarities {
		rank := i + 1
		if i > 0 {
			prevTokenRarity := tokenRarities[i-1]
			scoresEqual := IsFloat64Close(tokenRarity.Score(), prevTokenRarity.Score())
			if scoresEqual {
				if prevTokenRarity.Rank() == 0 {
					return nil, errors.New("preview token rarity rank is zero")
				}
				rank = prevTokenRarity.Rank()
			}
		}
		tokenRarity.SetRarityRanks(rank)
	}
	return tokenRarities, nil
}

// IsFloat64Close is used to judge whether two float64 are close enough
// It tries to be equivalent to the performance of math.isclose in python
// under the default parameters.
// todo: more precise implementation
func IsFloat64Close(a, b float64) bool {
	return (a == b) || (math.Abs(a-b) <= 1e-9)
}
