package handlers

import (
	"math"

	"github.com/Base-Labs/openrarity/models"
	"github.com/Base-Labs/openrarity/scoring"
)

// InformationContentScoringHandler implements the scoring.IScoreHandler.
// Rarity describes the information-theoretic "rarity" of a Collection.
// The concept of "rarity" can be considered as a measure of "surprise" at the
// occurrence of a particular token's properties, within the context of the
// Collection from which it is derived. Self-information is a measure of such
// surprise, and information entropy a measure of the expected value of
// self-information across a distribution (i.e. across a Collection).
//
// It is trivial to "stuff" a Collection with extra information by merely adding
// additional properties to all tokens. This is reflected in the Entropy field,
// measured in bits—all else held equal, a Collection with more token properties
// will have higher Entropy. However, this information bloat is carried by the
// tokens themselves, so their individual information-content grows in line with
// Collection-wide Entropy. The Scores are therefore scaled down by the Entropy
// to provide unless "relative surprise", which can be safely compared between
// Collections.
//
// This class computes rarity of each token in the Collection based on information
// entropy. Every TraitType is considered as a categorical probability
// distribution with each TraitValue having an associated probability and hence
// information content. The rarity of a particular token is the sum of
// information content carried by each of its Attributes, divided by the entropy
// of the Collection as a whole (see the Rarity struct for rationale).
//
// Notably, the lack of a TraitType is considered as a null-Value Attribute as
// the absence across the majority of a Collection implies rarity in those
// tokens that do carry the TraitType.
type InformationContentScoringHandler struct{}

var _ scoring.IScoreHandler = &InformationContentScoringHandler{}

// NewInformationContentScoringHandler is the constructor of InformationContentScoringHandler
func NewInformationContentScoringHandler() *InformationContentScoringHandler {
	return &InformationContentScoringHandler{}
}

// GetCollectionEntropy is used to Calculate the entropy of the collection,
// defined to be the sum of the probability of every possible attribute name/value
// pair that occurs in the collection times that square root of such probability.
func (c *InformationContentScoringHandler) GetCollectionEntropy(
	collection models.ICollection,
	attributes map[models.AttributeName][]*models.CollectionAttribute,
	nullAttributes map[models.AttributeName]*models.CollectionAttribute,
) float64 {
	if attributes == nil {
		attributes = collection.ExtractCollectionAttributes()
	}
	if nullAttributes == nil {
		nullAttributes = collection.ExtractNullAttributes()
	}
	collectionProbabilities := make([]float64, 0, len(attributes))
	for attrName, attrValues := range attributes {
		if nullAttr := nullAttributes[attrName]; nullAttr != nil {
			attrValues = append(attrValues, nullAttr)
		}
		probabilities := make([]float64, 0, len(attrValues))
		for _, attrValue := range attrValues {
			probabilities = append(probabilities,
				float64(attrValue.TotalTokens)/float64(collection.TokenTotalSupply()),
			)
		}
		collectionProbabilities = append(collectionProbabilities,
			probabilities...,
		)
	}
	var collectionEntropy float64
	for _, item := range collectionProbabilities {
		collectionEntropy += item * math.Log2(item)
	}
	return collectionEntropy * -1
}

// ScoreTokens should be used if you only want to score a batch of tokens that belong to collection.
// This will typically be more efficient than calling score_token for each
// token in `tokens`.
func (c *InformationContentScoringHandler) ScoreTokens(collection models.ICollection, tokens []models.IToken) ([]float64, error) {
	collectionNullAttributes := collection.ExtractNullAttributes()
	collectionAttributes := collection.ExtractCollectionAttributes()
	collectionEntropy := c.GetCollectionEntropy(
		collection,
		collectionAttributes,
		collectionNullAttributes,
	)
	if collectionEntropy == 0 {
		collectionEntropy = 1
	}
	scores := make([]float64, 0, len(tokens))
	for _, token := range tokens {
		scores = append(scores,
			c.scoreToken(
				collection,
				token,
				collectionNullAttributes,
				collectionEntropy,
			),
		)
	}
	return scores, nil
}

// ScoreToken is used to score an individual token based on the traits' distribution across
// the whole collection.
func (c *InformationContentScoringHandler) ScoreToken(
	collection models.ICollection,
	token models.IToken,
) (float64, error) {
	return c.scoreToken(
		collection, token,
		nil, 0,
	), nil
}

// scoreToken is used to calculate the score of the token using information
// entropy with a collection entropy normalization factor.
func (c *InformationContentScoringHandler) scoreToken(
	collection models.ICollection,
	token models.IToken,
	collectionNullAttributes map[models.AttributeName]*models.CollectionAttribute,
	collectionEntropyNormalization float64,
) float64 {
	icTokenScore := c.getICScore(collection, token, collectionNullAttributes)
	collectionEntropy := collectionEntropyNormalization
	if collectionEntropy == 0 {
		collectionEntropy = c.GetCollectionEntropy(
			collection,
			collection.ExtractCollectionAttributes(),
			collectionNullAttributes,
		)
	}
	normalizedTokenScore := icTokenScore / collectionEntropy
	return normalizedTokenScore
}

func (c *InformationContentScoringHandler) getICScore(
	collection models.ICollection,
	token models.IToken,
	collectionNullAttributes map[models.AttributeName]*models.CollectionAttribute,
) float64 {
	// First calculate the individual attribute scores for all attributes
	// of the provided token. Scores are the inverted probabilities of the
	// attribute in the collection.
	attrScores, _ := scoring.GetTokenAttributesScoresAndWeights(
		collection,
		token,
		false,
		collectionNullAttributes,
	)
	// Get a single score (via information content) for the token by taking
	// the sum of the logarithms of the attributes' scores.
	var result float64
	for _, score := range attrScores {
		result += math.Log2(1 / score)
	}
	return -1 * result
}
