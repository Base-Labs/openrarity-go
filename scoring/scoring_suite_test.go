package scoring_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/Base-Labs/openrarity"
	"github.com/Base-Labs/openrarity/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

// UniformRarityTokens returns a slice of IToken instances with uniform rarity.
// The number of attributes and the number of values per attribute are specified
// by the input arguments attributeCount and valuesPerAttribute, respectively.
// The total supply of tokens is determined by the tokenTotalSupply argument.
func UniformRarityTokens(
	attributeCount int,
	valuesPerAttribute int,
	tokenTotalSupply int,
) []models.IToken {
	tokens := make([]models.IToken, 0, tokenTotalSupply)
	for tokenId := 0; tokenId < tokenTotalSupply; tokenId++ {
		stringAttributeMap := map[string]models.IStringAttribute{}
		for i := 0; i < attributeCount; i++ {
			attrName := strconv.Itoa(i)
			stringAttributeMap[attrName] = models.NewStringAttribute(
				attrName,
				strconv.FormatInt(int64(tokenId/(tokenTotalSupply/valuesPerAttribute)), 10),
			)
		}
		tokens = append(tokens, models.NewToken(
			models.NewEVMContractTokenIdentifier("0x0", tokenId),
			models.TokenStandardERC721,
			models.NewTokenMetadataFromStringAttributes(stringAttributeMap),
		))
	}
	return tokens
}

// OneRareRarityTokens returns a slice of IToken instances with one rare token.
// The number of attributes and the number of values per attribute are specified
// by the input arguments attributeCount and valuesPerAttribute, respectively.
// The total supply of tokens is determined by the tokenTotalSupply argument.
func OneRareRarityTokens(
	attributeCount int,
	valuesPerAttribute int,
	tokenTotalSupply int,
) []models.IToken {
	tokens := make([]models.IToken, 0, tokenTotalSupply)
	for tokenId := 0; tokenId < tokenTotalSupply-1; tokenId++ {
		stringAttributeMap := map[string]models.IStringAttribute{}
		for i := 0; i < attributeCount; i++ {
			attrName := strconv.Itoa(i)
			stringAttributeMap[attrName] = models.NewStringAttribute(
				attrName,
				strconv.Itoa(
					tokenId/(tokenTotalSupply/(valuesPerAttribute-1))-1,
				),
			)
		}
		tokens = append(tokens, models.NewToken(
			models.NewEVMContractTokenIdentifier("0x0", tokenId),
			models.TokenStandardERC721,
			models.NewTokenMetadataFromStringAttributes(stringAttributeMap),
		))
	}
	rareTokenStringAttributeDict := make(map[models.AttributeName]models.IStringAttribute, attributeCount)
	for i := 0; i < attributeCount; i++ {
		attrName := strconv.Itoa(i)
		rareTokenStringAttributeDict[attrName] = models.NewStringAttribute(
			attrName,
			strconv.Itoa(valuesPerAttribute),
		)
	}
	tokens = append(tokens, models.NewToken(
		models.NewEVMContractTokenIdentifier("0x0", tokenTotalSupply-1),
		models.TokenStandardERC721,
		models.NewTokenMetadataFromStringAttributes(rareTokenStringAttributeDict),
	))
	return tokens
}

// Pair defines the struct of key-value pair
type Pair[Key any, Value any] struct {
	Key   Key
	Value Value
}

// GetMixedTraitSpread returns a piece of constructed data
func GetMixedTraitSpread(maxTotalSupply int) map[string][]Pair[string, float64] {
	totalSupply := float64(maxTotalSupply)
	return map[string][]Pair[string, float64]{
		"hat": {
			{"cap", float64(int(totalSupply * 0.2))},
			{"beanie", float64(int(totalSupply * 0.3))},
			{"hood", float64(int(totalSupply * 0.45))},
			{"visor", float64(int(totalSupply * 0.05))},
		},
		"shirt": {
			{"white-t", float64(int(totalSupply * 0.8))},
			{"vest", float64(int(totalSupply * 0.2))},
		},
		"special": {
			{"true", float64(int(totalSupply * 0.1))},
			{"null", float64(int(totalSupply * 0.9))},
		},
	}
}

// GenerateMixedCollection creates a new collection with a random mix of traits,
// using the given maximum total supply.
// The maximum total supply must be a multiple of 10 and greater than 100,
// otherwise this function will return error.
func GenerateMixedCollection(maxTotalSupply int) (models.ICollection, error) {
	if maxTotalSupply%10 != 0 || maxTotalSupply < 100 {
		return nil, errors.New("only multiples of 10 and greater than 100 please.")
	}
	tokenIDs := rand.Perm(maxTotalSupply)

	getTraitValue := func(traitSpread []Pair[string, float64], idx int) string {
		traitValueIdx := 0
		maxIdxForTraitValue := traitSpread[traitValueIdx].Value
		for float64(idx) >= maxIdxForTraitValue {
			traitValueIdx += 1
			maxIdxForTraitValue += traitSpread[traitValueIdx].Value
		}
		return traitSpread[traitValueIdx].Key
	}

	tokenIDsToTraits := make([]map[string]interface{}, maxTotalSupply)
	for idx, tokenID := range tokenIDs {
		traits := make(map[string]interface{})
		for traitName, traitValueToPercent := range GetMixedTraitSpread(10000) {
			traits[traitName] = getTraitValue(traitValueToPercent, idx)
		}
		tokenIDsToTraits[tokenID] = traits
	}

	return GenerateCollectionWithTokenTraits(tokenIDsToTraits, models.IdentifierTypeEVMContract)
}

// CreateEVMToken is used to create an evm token
func CreateEVMToken(
	tokenId int,
	contractAddress string,
	tokenStandard models.TokenStandard,
	metadata models.ITokenMetadata,
) openrarity.IToken {
	if metadata == nil {
		metadata = must(models.NewTokenMetadataFromAttributes(map[string]interface{}{}))
	}
	return openrarity.NewToken(
		models.NewEVMContractTokenIdentifier(
			contractAddress, tokenId,
		),
		tokenStandard,
		metadata,
	)
}

// GenerateCollectionWithTokenTraits is used to generate a new collection from the given traits
func GenerateCollectionWithTokenTraits(
	tokensTraits []map[string]interface{},
	tokenIdentifierType models.IdentifierType,
) (models.ICollection, error) {
	tokens := make([]models.IToken, 0, len(tokensTraits))
	for idx, tokenTraits := range tokensTraits {
		var (
			identifierType models.ITokenIdentifier
			tokenStandard  models.TokenStandard
		)
		switch tokenIdentifierType {
		case models.IdentifierTypeEVMContract:
			identifierType = models.NewEVMContractTokenIdentifier(
				"0x0", idx,
			)
			tokenStandard = models.TokenStandardERC721
		case models.IdentifierTypeSolanaMintAddress:
			identifierType = models.NewSolanaMintAddressTokenIdentifier(
				fmt.Sprintf("Fake-Address-%d", idx),
			)
			tokenStandard = models.TokenStandardERC721
		default:
			return nil, errors.Errorf("Unexpected token identifier type: %s", tokenIdentifierType)
		}
		tokenMetadata, err := models.NewTokenMetadataFromAttributes(tokenTraits)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, models.NewToken(
			identifierType,
			tokenStandard,
			tokenMetadata,
		))
	}
	return models.NewCollection("My Collection", tokens), nil
}

func must[V any](value V, err error) V {
	if err != nil {
		panic(err)
	}
	return value
}

func TestScoring(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scoring Suite")
}
