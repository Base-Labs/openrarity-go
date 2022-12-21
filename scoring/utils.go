package scoring

import (
	"sort"

	"github.com/Base-Labs/openrarity/models"
)

// GetTokenAttributesScoresAndWeights is used to calculate the scores and normalization weights for a token
// based on its attributes. If the token does not have an attribute, the probability
// of the attribute being null is used instead.
func GetTokenAttributesScoresAndWeights(
	collection models.ICollection,
	token models.IToken,
	normalized bool,
	collectionNullAttributes map[models.AttributeName]*models.CollectionAttribute,
) ([]float64, []float64) {
	nullAttributes := collectionNullAttributes
	if nullAttributes == nil {
		nullAttributes = collection.ExtractNullAttributes()
	}
	combinedAttributes := MergeMap(nullAttributes, convertToCollectionAttributesDict(collection, token))
	sortedAttrNames := GetMapKeys(combinedAttributes)
	sort.Slice(sortedAttrNames, func(i, j int) bool {
		return sortedAttrNames[i] < sortedAttrNames[j]
	})
	sortedAttrs := make([]*models.CollectionAttribute, 0, len(sortedAttrNames))
	for _, name := range sortedAttrNames {
		sortedAttrs = append(sortedAttrs, combinedAttributes[name])
	}
	totalSupply := collection.TokenTotalSupply()

	attrWeights := make([]float64, 0, len(sortedAttrNames))
	if normalized {
		for _, attrName := range sortedAttrNames {
			attrWeights = append(attrWeights, float64(1)/float64(collection.TotalAttributeValues(attrName)))
		}
	} else {
		for range sortedAttrNames {
			attrWeights = append(attrWeights, float64(len(sortedAttrNames)))
		}
	}
	scores := make([]float64, 0, len(sortedAttrs))
	for _, attr := range sortedAttrs {
		scores = append(scores, float64(totalSupply)/float64(attr.TotalTokens))
	}
	return scores, attrWeights
}

// GetMapKeys is used to all keys in a map
func GetMapKeys[K comparable, V any](a map[K]V) []K {
	data := make([]K, 0, len(a))
	for k := range a {
		data = append(data, k)
	}
	return data
}

// MergeMap is used to merge map
func MergeMap[K comparable, V any](a map[K]V, b map[K]V) map[K]V {
	data := make(map[K]V, (len(a)+len(b))*2/3)
	for k, v := range a {
		data[k] = v
	}
	for k, v := range b {
		data[k] = v
	}
	return data
}

func convertToCollectionAttributesDict(collection models.ICollection, token models.IToken) map[models.AttributeName]*models.CollectionAttribute {
	// We currently only support string attributes
	attributes := make(map[models.AttributeName]*models.CollectionAttribute, len(token.Metadata().StringAttributes()))
	for name, value := range token.Metadata().StringAttributes() {
		attributes[name] = &models.CollectionAttribute{
			Attribute:   value,
			TotalTokens: collection.TotalTokensWithAttributes(value),
		}
	}
	return attributes
}
