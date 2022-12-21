package models

// ITokenRarity holds rarity and optional rank information along with the token
type ITokenRarity interface {
	// Score is used to obtain the rarity score of the current token.
	Score() float64
	// TokenFeatures is used to get various features of current token.
	TokenFeatures() ITokenRankingFeatures
	// SetRarityRanks is used to set the rarity ranking of the current token.
	SetRarityRanks(rank int)
	// Token is used to return the current token.
	Token() IToken
	// Rank is used to obtain the rarity ranking of the current token.
	Rank() int
}

// TokenRarity hold rarity and optional rank information along with the token
type TokenRarity struct {
	score         float64
	tokenFeatures ITokenRankingFeatures
	token         IToken
	rank          int
}

var _ ITokenRarity = &TokenRarity{}

// NewTokenRarity is the constructor of TokenRarity
func NewTokenRarity(
	token IToken,
	score float64,
	tokenFeatures ITokenRankingFeatures,
) *TokenRarity {
	return &TokenRarity{
		score:         score,
		tokenFeatures: tokenFeatures,
		token:         token,
	}
}

// SetRarityRanks is used to set the rarity ranking of the current token.
func (c *TokenRarity) SetRarityRanks(rank int) {
	c.rank = rank
}

// Score is used to obtain the rarity score of the current token.
func (c *TokenRarity) Score() float64 {
	return c.score
}

// TokenFeatures is used to get various features of current token.
func (c *TokenRarity) TokenFeatures() ITokenRankingFeatures {
	return c.tokenFeatures
}

// Token is used to return the current token.
func (c *TokenRarity) Token() IToken {
	return c.token
}

// Rank is used to obtain the rarity ranking of the current token.
func (c *TokenRarity) Rank() int {
	return c.rank
}
