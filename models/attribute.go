package models

import (
	"strings"
)

// TraitCountAttributeName defines the trait count attribute name
const TraitCountAttributeName = "meta_trait:trait_count"

// IAttribute defines the interface of attribute
type IAttribute interface {
	// Name returns name of an attribute
	Name() AttributeName
}

// IStringAttribute represent string token attribute name and value
type IStringAttribute interface {
	// Name returns name of an attribute
	Name() AttributeName
	// Value returns value of a string attribute
	Value() StringAttributeValue
}

// INumericAttribute represent numeric token attribute name and value
type INumericAttribute interface {
	// Name returns name of an attribute
	Name() AttributeName
	// Value returns value of a numeric attribute
	Value() INumericAttributeValue
}

// IDateAttribute represent date token attribute name and value
type IDateAttribute interface {
	// Name returns name of an attribute
	Name() AttributeName
	// Value returns value of a date attribute
	Value() DateAttributeValue
}

// AttributeName defines the type of attribute name
type AttributeName = string

// NormalizeAttributeString is used to normalize the key-value of the attribute
func NormalizeAttributeString(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

// StringAttributeValue defines the type of attribute value
type StringAttributeValue = string

// StringAttribute represent string token attribute name and value
type StringAttribute struct {
	name  AttributeName
	value StringAttributeValue
}

var _ IStringAttribute = &StringAttribute{}
var _ IAttribute = &StringAttribute{}

// NewStringAttribute is the constructor of StringAttribute
func NewStringAttribute(name string, value string) StringAttribute {
	return StringAttribute{
		name:  AttributeName(NormalizeAttributeString(name)),
		value: StringAttributeValue(NormalizeAttributeString(value)),
	}
}

// Name returns name of an attribute
func (s StringAttribute) Name() AttributeName {
	return s.name
}

// Value returns value of a string attribute
func (s StringAttribute) Value() StringAttributeValue {
	return s.value
}

// NumericAttribute represent numeric token attribute name and value
type NumericAttribute struct {
	name  AttributeName
	value INumericAttributeValue
}

var _ INumericAttribute = &NumericAttribute{}

// NewNumericAttribute is the constructor of NumericAttribute
func NewNumericAttribute[V float64 | int64 | int](name string, value V) *NumericAttribute {
	return &NumericAttribute{
		name:  AttributeName(NormalizeAttributeString(name)),
		value: NewNumericAttributeValue(value),
	}
}

// Name returns name of an attribute
func (c NumericAttribute) Name() AttributeName {
	return c.name
}

// Value returns value of a numeric attribute
func (c NumericAttribute) Value() INumericAttributeValue {
	return c.value
}

// INumericAttributeValue represent numeric token attribute value
type INumericAttributeValue interface {
	// Float64 is used to get the stored data.
	// If the stored data is of floating point type, the original value can be obtained by this method
	Float64() (float64, bool)
	// Int64 is used to get the stored data.
	// If the stored data is of integer type, the original value can be obtained by this method
	Int64() (int64, bool)
}

var _ INumericAttributeValue = &NumericAttributeValue{}

// NumericAttributeValue represent numeric token attribute value
type NumericAttributeValue struct {
	value interface{}
}

// NewNumericAttributeValue creates a new NumericAttributeValue object that can hold
// a value of type float64, integer.
func NewNumericAttributeValue[V float64 | int64 | int](val V) *NumericAttributeValue {
	return &NumericAttributeValue{
		value: val,
	}
}

// Float64 is used to get the stored data.
// If the stored data is of floating point type, the original value can be obtained by this method
func (n NumericAttributeValue) Float64() (float64, bool) {
	switch v := n.value.(type) {
	case float64:
		return v, true
	default:
		return 0, false
	}
}

// Int64 is used to get the stored data.
// If the stored data is of integer type, the original value can be obtained by this method
func (n NumericAttributeValue) Int64() (int64, bool) {
	switch v := n.value.(type) {
	case int64:
		return v, true
	case int:
		return int64(v), true
	default:
		return 0, false
	}
}

// DateAttribute represent date token attribute name and value
type DateAttribute struct {
	name  AttributeName
	value DateAttributeValue
}

// DateAttributeValue defines the value type of DateAttribute
type DateAttributeValue = int64

var _ IDateAttribute = &DateAttribute{}

// NewDateAttribute is the constructor of NewDateAttribute
func NewDateAttribute(name string, value int64) *DateAttribute {
	return &DateAttribute{
		name:  AttributeName(NormalizeAttributeString(name)),
		value: DateAttributeValue(value),
	}
}

// Name returns name of an attribute
func (c *DateAttribute) Name() AttributeName {
	return c.name
}

// Value returns value of a date attribute
func (c *DateAttribute) Value() DateAttributeValue {
	return c.value
}
