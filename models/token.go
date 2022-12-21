package models

// IToken represents a token on the blockchain.
// Examples of these are non-fungible tokens, or semi-fungible tokens.
type IToken interface {
	// Metadata is used to get the metadata of current token.
	Metadata() ITokenMetadata
	// TraitCount is used to return the count of non-null, non-"none" value traits this token has.
	TraitCount() int
	// HasAttribute is used to determine whether the metadata of the current Token
	// contains the given attribute name.
	HasAttribute(name AttributeName) bool
	// TokenStandard is used to obtain the current Token standard.
	TokenStandard() TokenStandard
	// TokenIdentifier is used to obtain the current Token identifier.
	TokenIdentifier() ITokenIdentifier
}

var _ IToken = &Token{}

// Token represents a token on the blockchain.
// Examples of these are non-fungible tokens, or semi-fungible tokens.
type Token struct {
	metadata        ITokenMetadata
	tokenStandard   TokenStandard
	tokenIdentifier ITokenIdentifier
}

// NewToken is the constructor of Token
func NewToken(
	tokenIdentifier ITokenIdentifier,
	tokenStandard TokenStandard,
	metadata ITokenMetadata,
) *Token {
	return &Token{
		tokenIdentifier: tokenIdentifier,
		tokenStandard:   tokenStandard,
		metadata:        metadata,
	}
}

// NewERC721Token Creates a Token class representing an ERC721 evm token given the following
// parameters.
func NewERC721Token(
	contractAddress string,
	tokenID int,
	metadata map[string]interface{},
) (*Token, error) {
	attributes, err := NewTokenMetadataFromAttributes(metadata)
	if err != nil {
		return nil, err
	}
	return &Token{
		tokenIdentifier: NewEVMContractTokenIdentifier(
			contractAddress,
			tokenID,
		),
		tokenStandard: TokenStandardERC721,
		metadata:      attributes,
	}, nil
}

// HasAttribute is used to determine whether the metadata of the current Token
// contains the given attribute name.
func (c *Token) HasAttribute(name AttributeName) bool {
	return c.metadata.AttributeExists(name)
}

// TokenStandard is used to obtain the current Token standard.
func (c *Token) TokenStandard() TokenStandard {
	return c.tokenStandard
}

// Metadata is used to get the metadata of current token.
func (c *Token) Metadata() ITokenMetadata {
	return c.metadata
}

// TraitCount is used to return the count of non-null, non-"none" value traits this token has.
func (c *Token) TraitCount() int {
	return GetStringAttributesCount(c.metadata.StringAttributes()) +
		len(c.metadata.NumericAttributes()) +
		len(c.metadata.DateAttributes())
}

// TokenIdentifier is used to obtain the current Token identifier.
func (c *Token) TokenIdentifier() ITokenIdentifier {
	return c.tokenIdentifier
}

// GetStringAttributesCount returns the number of string attributes in the given map
// that have a non-null and non-"none" value after normalization.
func GetStringAttributesCount(attributes map[AttributeName]IStringAttribute) int {
	var count int
	for _, attribute := range attributes {
		switch NormalizeAttributeString(attribute.Value()) {
		case "", "none":
		default:
			count++
		}
	}
	return count
}
