package scoring_test

import (
	"fmt"
	"math"
	"time"

	"github.com/Base-Labs/openrarity/models"
	"github.com/Base-Labs/openrarity/scoring"
	"github.com/Base-Labs/openrarity/scoring/handlers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scoring Handlers", func() {
	maxScoringTimeFor10ks := 2

	uniformTokens := UniformRarityTokens(10000, 5, 10)
	uniformCollection := models.NewCollection("", uniformTokens)

	oneRareTokens := OneRareRarityTokens(10000, 3, 10)
	oneRareCollection := models.NewCollection("", oneRareTokens)

	mixedCollection, err := GenerateMixedCollection(10000)
	Expect(err).To(BeNil())

	It("should able to pass test_information_content_rarity_uniform", func() {
		icHandler := handlers.NewInformationContentScoringHandler()
		uniformTokenToTest := uniformCollection.Tokens()[0]
		uniformIcRarity := 1
		score0, err := icHandler.ScoreToken(uniformCollection, uniformTokenToTest)
		Expect(err).To(BeNil())
		Expect(int(math.Round(score0*1e8) / 1e8)).To(Equal(uniformIcRarity))
	})

	It("should able to pass test_information_content_null_value_attribute", func() {
		icScorer := handlers.NewInformationContentScoringHandler()
		collectionWithEmpty, err := GenerateCollectionWithTokenTraits(
			[]map[string]interface{}{
				{"bottom": "spec", "hat": "spec", "special": "true"},
				{"bottom": "1", "hat": "1", "special": "true"},
				{"bottom": "1", "hat": "1"},
				{"bottom": "2", "hat": "2"},
				{"bottom": "2", "hat": "2"},
				{"bottom": "3", "hat": "2"},
			},
			models.IdentifierTypeEVMContract,
		)
		Expect(err).To(BeNil())

		collectionEntropy := icScorer.GetCollectionEntropy(collectionWithEmpty, nil, nil)
		var collectionProbs []float64
		spread := map[string]map[string]int{
			"bottom":                       {"1": 2, "2": 2, "3": 1, "spec": 1},
			"hat":                          {"1": 2, "2": 3, "spec": 1},
			"special":                      {"true": 2, "Null": 4},
			models.TraitCountAttributeName: {"2": 4, "3": 2},
		}
		for _, traitDict := range spread {
			for _, tokensWithTrait := range traitDict {
				collectionProbs = append(collectionProbs, float64(tokensWithTrait)/float64(6))
			}
		}
		var expectedCollectionEntropy float64
		for _, item := range collectionProbs {
			expectedCollectionEntropy += item * math.Log2(item)
		}
		expectedCollectionEntropy = expectedCollectionEntropy * -1
		Expect(fmt.Sprintf("%.10f", expectedCollectionEntropy)).To(
			Equal(fmt.Sprintf("%.10f", collectionEntropy)))

		scores, err := icScorer.ScoreTokens(
			collectionWithEmpty,
			collectionWithEmpty.Tokens(),
		)
		Expect(err).To(BeNil())

		Expect(scores[0] > scores[1]).To(BeTrue())
		Expect(scores[1] > scores[2]).To(BeTrue())
		Expect(scores[5] > scores[2]).To(BeTrue())
		Expect(scores[2] > scores[3]).To(BeTrue())
		Expect(scores[3] == scores[4]).To(BeTrue())

		for i, token := range collectionWithEmpty.Tokens() {
			attrScores, _ := scoring.GetTokenAttributesScoresAndWeights(
				collectionWithEmpty,
				token,
				false,
				nil,
			)
			var icTokenScore float64
			for _, score := range attrScores {
				icTokenScore += math.Log2(1 / score)
			}
			icTokenScore = -1 * icTokenScore
			expectedScore := icTokenScore / collectionEntropy
			Expect(fmt.Sprintf("%.10f", scores[i])).To(
				Equal(fmt.Sprintf("%.10f", expectedScore)))
		}
	})

	It("should able to pass test_information_content_empty_attribute", func() {
		collectionWithNull, err := GenerateCollectionWithTokenTraits(
			[]map[string]interface{}{
				{"bottom": "1", "hat": "1", "special": "true"},
				{"bottom": "1", "hat": "1"},
				{"bottom": "2", "hat": "2"},
				{"bottom": "2", "hat": "2"},
				{"bottom": "3", "hat": "2"},
			},
			models.IdentifierTypeEVMContract,
		)
		Expect(err).To(BeNil())

		collectionWithoutNull, err := GenerateCollectionWithTokenTraits(
			[]map[string]interface{}{
				{"bottom": "1", "hat": "1", "special": "true"},
				{"bottom": "1", "hat": "1", "special": "none"},
				{"bottom": "2", "hat": "2", "special": "none"},
				{"bottom": "2", "hat": "2", "special": "none"},
				{"bottom": "3", "hat": "2", "special": "none"},
			},
			models.IdentifierTypeEVMContract,
		)
		Expect(err).To(BeNil())

		icScorer := handlers.NewInformationContentScoringHandler()
		scoresWithNull, err := icScorer.ScoreTokens(
			collectionWithNull,
			collectionWithNull.Tokens(),
		)
		Expect(err).To(BeNil())

		scoresWithoutNull, err := icScorer.ScoreTokens(
			collectionWithoutNull,
			collectionWithoutNull.Tokens(),
		)
		Expect(err).To(BeNil())

		for idx, value := range scoresWithNull {
			Expect(fmt.Sprintf("%.10f", value)).To(
				Equal(fmt.Sprintf("%.10f", scoresWithoutNull[idx])))
		}
	})
	It("should able to pass test_information_content_rarity_timing", func() {
		icScorer := handlers.NewInformationContentScoringHandler()
		tic := time.Now()
		_, err := icScorer.ScoreTokens(
			mixedCollection,
			mixedCollection.Tokens(),
		)
		Expect(err).To(BeNil())
		toc := time.Now()
		Expect(toc.Sub(tic).Seconds() < float64(maxScoringTimeFor10ks)).To(BeTrue())
	})

	_ = oneRareCollection
})
