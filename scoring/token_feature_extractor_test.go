package scoring_test

import (
	"github.com/Base-Labs/openrarity/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Feature Extractor", func() {
	It("should pass test_feature_extractor", func() {
		collection, err := GenerateCollectionWithTokenTraits(
			[]map[string]interface{}{
				{"bottom": "1", "hat": "1", "special": "true"},
				{"bottom": "1", "hat": "1", "special": "false"},
				{"bottom": "2", "hat": "2", "special": "false"},
				{"bottom": "2", "hat": "2", "special": "false"},
				{"bottom": "3", "hat": "2", "special": "false"},
				{"bottom": "4", "hat": "3", "special": "false"},
			},
			models.IdentifierTypeEVMContract,
		)
		Expect(err).To(BeNil())

		expectedCount := []int{1, 0, 0, 0, 1, 2}
		for i := 0; i < 5; i++ {
			Expect(models.ExtractUniqueAttributeCount(
				collection.Tokens()[i], collection,
			).UniqueAttributeCount()).To(Equal(expectedCount[i]))
		}
	})
})
