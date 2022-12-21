package models

import (
	"time"

	"github.com/pkg/errors"
)

// ITokenMetadata represent EIP-721 or EIP-1115 compatible metadata structure
type ITokenMetadata interface {
	// StringAttributes is returns the mapping of attribute name string attribute value
	StringAttributes() map[AttributeName]IStringAttribute
	// AddAttribute is used to add an attribute to this metadata object, overriding existing
	// attribute if the normalized attribute name already exists.
	AddAttribute(attribute IAttribute)
	// AttributeExists returns True if this metadata object has an attribute with the given name.
	AttributeExists(name AttributeName) bool
	// NumericAttributes is returns the mapping of attribute name numeric attribute value
	NumericAttributes() map[AttributeName]INumericAttribute
	// DateAttributes is returns the mapping of attribute name date attribute value
	DateAttributes() map[AttributeName]IDateAttribute
}

var _ ITokenMetadata = &TokenMetadata{}

// TokenMetadata represent EIP-721 or EIP-1115 compatible metadata structure
type TokenMetadata struct {
	stringAttributes  map[AttributeName]IStringAttribute
	numericAttributes map[AttributeName]INumericAttribute
	dateAttributes    map[AttributeName]IDateAttribute
}

// NewTokenMetadataFromStringAttributes is used to create string attributes from stringAttributes
func NewTokenMetadataFromStringAttributes(
	stringAttributes map[AttributeName]IStringAttribute) *TokenMetadata {
	return &TokenMetadata{
		stringAttributes: stringAttributes,
	}
}

// NewTokenMetadataFromAttributes is used to create TokenMetadata from attributes
func NewTokenMetadataFromAttributes(attributes map[string]interface{}) (*TokenMetadata, error) {
	stringAttributes := map[AttributeName]IStringAttribute{}
	numericAttributes := map[AttributeName]INumericAttribute{}
	dateAttributes := map[AttributeName]IDateAttribute{}
	for attrName, attrValue := range attributes {
		normalizeAttributeName := NormalizeAttributeString(attrName)
		switch v := attrValue.(type) {
		case string:
			stringAttributes[normalizeAttributeName] = NewStringAttribute(
				normalizeAttributeName, v,
			)
		case float64:
			numericAttributes[normalizeAttributeName] = NewNumericAttribute(
				normalizeAttributeName, v,
			)
		case int64:
			numericAttributes[normalizeAttributeName] = NewNumericAttribute(
				normalizeAttributeName, v,
			)
		case int:
			numericAttributes[normalizeAttributeName] = NewNumericAttribute(
				normalizeAttributeName, v,
			)
		case time.Time:
			dateAttributes[normalizeAttributeName] = NewDateAttribute(
				normalizeAttributeName, v.Unix(),
			)
		default:
			return nil, errors.Errorf("Provided attribute value has invalid type: %T, Must be string.", v)
		}
	}
	return &TokenMetadata{
		stringAttributes:  stringAttributes,
		numericAttributes: numericAttributes,
		dateAttributes:    dateAttributes,
	}, nil
}

// NumericAttributes is returns the mapping of attribute name numeric attribute value
func (c *TokenMetadata) NumericAttributes() map[AttributeName]INumericAttribute {
	return c.numericAttributes
}

// DateAttributes is returns the mapping of attribute name date attribute value
func (c *TokenMetadata) DateAttributes() map[AttributeName]IDateAttribute {
	return c.dateAttributes
}

// StringAttributes is returns tje mapping of attribute name to list of string attribute values
func (c *TokenMetadata) StringAttributes() map[AttributeName]IStringAttribute {
	return c.stringAttributes
}

// AttributeExists returns True if this metadata object has an attribute with the given name.
func (c *TokenMetadata) AttributeExists(name AttributeName) bool {
	if _, exists := c.stringAttributes[name]; exists {
		return true
	}
	return false
}

// AddAttribute is used to add an attribute to this metadata object, overriding existing
// attribute if the normalized attribute name already exists.
func (c *TokenMetadata) AddAttribute(attribute IAttribute) {
	switch v := attribute.(type) {
	case IStringAttribute:
		c.stringAttributes[attribute.Name()] = v
	}
}
