package main

import (
	"fmt"

	"github.com/Base-Labs/openrarity"
)

func must[V any](value V, err error) V {
	if err != nil {
		panic(err)
	}
	return value
}

func main() {
	scorer := openrarity.NewOpenRarityScorer()
	collection := openrarity.NewCollection(
		"My Collection Name",
		[]openrarity.IToken{
			must(openrarity.NewERC721Token("0xa3049...", 1, map[string]interface{}{
				"hat":   "cap",
				"shirt": "blue",
			})),
			must(openrarity.NewERC721Token("0xa3049...", 2, map[string]interface{}{
				"hat":   "visor",
				"shirt": "green",
			})),
			must(openrarity.NewERC721Token("0xa3049...", 3, map[string]interface{}{
				"hat":   "visor",
				"shirt": "blue",
				"color": "blue",
			})),
		},
	)
	ranker := openrarity.NewRarityRanker()
	rankedTokens, err := ranker.RankCollection(collection, scorer)
	if err != nil {
		panic(err)
	}
	for _, token := range rankedTokens {
		fmt.Println(token.Token().TokenIdentifier(), token.Rank(), token.Score())
	}
}
