package models

import (
	"strconv"
)

// ICollection represents collection of tokens used to determine token rarity score.
// A token's rarity is influenced by the attribute frequency of all the tokens
// in a collection.
type ICollection interface {
	// Tokens method is used to get all tokens in this collection
	Tokens() []IToken
	// TokenTotalSupply is used get the total supply of this collection
	TokenTotalSupply() int
	// TotalAttributeValues is used to get the number of values of specified attributeName
	TotalAttributeValues(attributeName AttributeName) int
	// TotalTokensWithAttributes is used to return the numbers of tokens in this collection with the attribute
	// based on the attributes' frequency counts.
	TotalTokensWithAttributes(attribute IStringAttribute) int
	// ExtractNullAttributes is used to compute probabilities of Null attributes.
	ExtractNullAttributes() map[AttributeName]*CollectionAttribute
	// ExtractCollectionAttributes is used to extract the map of collection traits with its respective counts
	ExtractCollectionAttributes() map[AttributeName][]*CollectionAttribute
	// TokenStandards is used to return token standards for this collection.
	TokenStandards() []TokenStandard
	// HasNumericAttribute is used to determine whether the current collection contains numeric attributes
	HasNumericAttribute() bool
}

var _ ICollection = &Collection{}

// Collection represents collection of tokens used to determine token rarity score.
// A token's rarity is influenced by the attribute frequency of all the tokens
// in a collection.
type Collection struct {
	name                      string
	tokens                    []IToken
	attributesFrequencyCounts map[AttributeName]map[StringAttributeValue]int
}

// CollectionAttribute represents an attribute that at least one token in a Collection has.
// E.g. "hat" = "cap" would be one attribute, and "hat" = "beanie" would be another
// unique attribute, even though they may belong to the same attribute type (id=name).
type CollectionAttribute struct {
	Attribute   IStringAttribute
	TotalTokens int
}

// NewCollection is the constructor of Collection
func NewCollection(name string, tokens []IToken) *Collection {
	c := &Collection{
		name:   name,
		tokens: tokens,
	}
	c.traitCountify(tokens)
	c.attributesFrequencyCounts = c.deriveNormalizedAttrsFrequencyCount()
	return c
}

// HasNumericAttribute is used to determine whether the current collection contains numeric attributes
func (c *Collection) HasNumericAttribute() bool {
	for _, token := range c.tokens {
		if len(token.Metadata().NumericAttributes()) > 0 ||
			len(token.Metadata().DateAttributes()) > 0 {
			return true
		}
	}
	return false
}

// TotalTokensWithAttributes is used to return the numbers of tokens in this collection with the attribute
// based on the attributes' frequency counts.
func (c *Collection) TotalTokensWithAttributes(attribute IStringAttribute) int {
	return c.attributesFrequencyCounts[attribute.Name()][attribute.Value()]
}

// deriveNormalizedAttrsFrequencyCount is used to Derive and construct attributes_frequency_counts based on
// string attributes on tokens. Numeric or date attributes currently not
// supported.
func (c *Collection) deriveNormalizedAttrsFrequencyCount() map[AttributeName]map[StringAttributeValue]int {
	attrsFreqCounts := map[AttributeName]map[StringAttributeValue]int{}
	for _, token := range c.tokens {
		for attrName, strAttr := range token.Metadata().StringAttributes() {
			if attrsFreqCounts[attrName] == nil {
				attrsFreqCounts[attrName] = map[StringAttributeValue]int{
					strAttr.Value(): 1,
				}
			} else {
				attrsFreqCounts[attrName][strAttr.Value()]++
			}
		}
	}
	return attrsFreqCounts
}

// traitCountify is used to Update tokens to have meta attribute "meta trait: trait_count" if it doesn't
// already exist.
func (c *Collection) traitCountify(tokens []IToken) {
	for _, token := range tokens {
		traitCount := token.TraitCount()
		if token.HasAttribute(TraitCountAttributeName) {
			traitCount--
		}
		token.Metadata().AddAttribute(
			NewStringAttribute(
				TraitCountAttributeName,
				strconv.FormatInt(int64(traitCount), 10),
			),
		)
	}
}

// Tokens method is used to get all tokens in this collection
func (c *Collection) Tokens() []IToken {
	return c.tokens
}

// TokenTotalSupply is used get the total supply of this collection
func (c *Collection) TokenTotalSupply() int {
	return len(c.tokens)
}

// TotalAttributeValues is used to get the number of values of specified attributeName
func (c *Collection) TotalAttributeValues(attributeName AttributeName) int {
	return len(c.attributesFrequencyCounts[attributeName])
}

// ExtractNullAttributes is used to compute probabilities of Null attributes.
func (c *Collection) ExtractNullAttributes() map[AttributeName]*CollectionAttribute {
	result := map[AttributeName]*CollectionAttribute{}
	for traitName, traitValues := range c.attributesFrequencyCounts {
		var totalTraitCount int
		for _, count := range traitValues {
			totalTraitCount += count
		}
		assetsWithoutTrait := c.TokenTotalSupply() - totalTraitCount
		if assetsWithoutTrait > 0 {
			result[traitName] = &CollectionAttribute{
				Attribute:   NewStringAttribute(traitName, "Null"),
				TotalTokens: assetsWithoutTrait,
			}
		}
	}
	return result
}

// ExtractCollectionAttributes is used to extract the map of collection traits with its respective counts
func (c *Collection) ExtractCollectionAttributes() map[AttributeName][]*CollectionAttribute {
	collectionTraits := map[AttributeName][]*CollectionAttribute{}
	for traitName, traitValues := range c.attributesFrequencyCounts {
		for traitValue, traitCount := range traitValues {
			collectionTraits[traitName] = append(
				collectionTraits[traitName],
				&CollectionAttribute{
					Attribute:   NewStringAttribute(traitName, string(traitValue)),
					TotalTokens: traitCount,
				},
			)
		}
	}
	return collectionTraits
}

// TokenStandards is used to return token standards for this collection.
func (c *Collection) TokenStandards() []TokenStandard {
	tokenStandards := NewSet[TokenStandard](len(c.Tokens()))
	for _, token := range c.Tokens() {
		tokenStandards.Add(token.TokenStandard())
	}
	return tokenStandards.List()
}
