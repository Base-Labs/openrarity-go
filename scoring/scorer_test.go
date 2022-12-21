package scoring_test

import (
	"github.com/Base-Labs/openrarity"
	"github.com/Base-Labs/openrarity/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scoring", func() {
	It("should pass test_score_collections_with_string_attributes", func() {
		collection, err := GenerateCollectionWithTokenTraits(
			[]map[string]interface{}{
				{"bottom": "1", "hat": "1", "special": "true"},
				{"bottom": "1", "hat": "1", "special": "false"},
				{"bottom": "2", "hat": "2", "special": "false"},
				{"bottom": "2", "hat": "2", "special": "false"},
				{"bottom": "pants", "hat": "cap", "special": "false"},
			},
			models.IdentifierTypeEVMContract,
		)
		Expect(err).To(BeNil())

		collectionTwo, err := GenerateCollectionWithTokenTraits(
			[]map[string]interface{}{
				{"bottom": "1", "hat": "1", "special": "true"},
				{"bottom": "1", "hat": "1", "special": "false"},
				{"bottom": "2", "hat": "2", "special": "false"},
				{"bottom": "2", "hat": "2", "special": "false"},
				{"bottom": "pants", "hat": "cap", "special": "false"},
			},
			models.IdentifierTypeEVMContract,
		)
		Expect(err).To(BeNil())

		scorer := openrarity.NewOpenRarityScorer()

		scores, err := scorer.ScoreCollection(collection)
		Expect(err).To(BeNil())
		Expect(len(scores)).To(Equal(5))

		scores, err = scorer.ScoreTokens(collection, collection.Tokens())
		Expect(err).To(BeNil())
		Expect(len(scores)).To(Equal(5))

		scoresTwo, err := scorer.ScoreCollections([]models.ICollection{
			collection, collectionTwo,
		})
		Expect(len(scoresTwo[0])).To(Equal(5))
		Expect(len(scoresTwo[1])).To(Equal(5))
	})
	It("should pass test_score_collection_with_numeric_attribute_errors", func() {
		collection, err := GenerateCollectionWithTokenTraits(
			[]map[string]interface{}{
				{"bottom": "1", "hat": "1", "special": "true"},
				{"bottom": "1", "hat": "1", "special": "false"},
				{"bottom": "2", "hat": "2", "special": "false"},
				{"bottom": "2", "hat": "2", "special": "false"},
				{"bottom": 3, "hat": 2, "special": "false"},
			},
			models.IdentifierTypeEVMContract,
		)
		Expect(err).To(BeNil())

		scorer := openrarity.NewOpenRarityScorer()
		scores, err := scorer.ScoreCollection(collection)
		Expect(scores).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("OpenRarity currently does not support collections " +
			"with numeric or date traits"))
	})
	It("should pass test_score_collection_with_erc1155_errors", func() {
		tokens := make([]openrarity.IToken, 0, 10)
		for i := 0; i < 10; i++ {
			tokens = append(tokens, CreateEVMToken(
				i,
				"0xaaa",
				models.TokenStandardERC1155,
				nil,
			))
		}
		collection := openrarity.NewCollection("test", tokens)
		scorer := openrarity.NewOpenRarityScorer()
		scores, err := scorer.ScoreCollection(collection)
		Expect(scores).To(BeNil())
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("OpenRarity currently only supports " +
			"ERC721/Non-fungible standards"))
	})
})
