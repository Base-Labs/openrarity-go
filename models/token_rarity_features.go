package models

// ITokenRankingFeatures is used to extract features from tokens
type ITokenRankingFeatures interface {
	// UniqueAttributeCount is used to get the unique attribute count
	UniqueAttributeCount() int
}

// TokenRankingFeatures is used to extract features from tokens
type TokenRankingFeatures struct {
	uniqueAttributeCount int
}

// NewTokenRankingFeatures is the constructor of TokenRankingFeatures
func NewTokenRankingFeatures(uniqueAttributeCount int) *TokenRankingFeatures {
	return &TokenRankingFeatures{
		uniqueAttributeCount: uniqueAttributeCount,
	}
}

var _ ITokenRankingFeatures = &TokenRankingFeatures{}

// UniqueAttributeCount is used to get the unique attribute count
func (c *TokenRankingFeatures) UniqueAttributeCount() int {
	return c.uniqueAttributeCount
}

// ExtractUniqueAttributeCount is used to extract unique attributes count from the token
func ExtractUniqueAttributeCount(
	token IToken, collection ICollection,
) ITokenRankingFeatures {
	uniqueAttributesCount := 0
	for _, stringAttribute := range token.Metadata().StringAttributes() {
		count := collection.TotalTokensWithAttributes(stringAttribute)
		if count == 1 {
			uniqueAttributesCount++
		}
	}
	return NewTokenRankingFeatures(uniqueAttributesCount)
}
