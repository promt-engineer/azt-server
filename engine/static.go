package engine

import (
	"aztec-pyramids/engine/models"
	"github.com/samber/lo"
	"math/rand"
)

const (
	extraSpinCount                  = 5
	needScattersToTriggerBonus      = 4
	needScattersToTriggerExtraSpins = 3

	topReelHeight = 4
	windowWidth   = 6

	reelSet1Code = 0
	reelSet2Code = 1
	reelSet3Code = 2

	topReelIndex = 6

	multiplicationDivider = 10

	buyBonusMultiplier = 100

	anteBetMultiplier = 125
	anteBetDivider    = 100

	buyBonusScatterQty = 4

	avgFreeWinBuyBonus = 96
)

var wildSymbol = 11
var scatterSymbol = 12

var multipliers = map[int]map[int]int64{
	1:  {2: 10, 3: 20, 4: 100, 5: 250, 6: 500},
	2:  {3: 10, 4: 20, 5: 25, 6: 50},
	3:  {3: 3, 4: 5, 5: 10, 6: 25},
	4:  {3: 3, 4: 5, 5: 8, 6: 20},
	5:  {3: 2, 4: 4, 5: 6, 6: 15},
	6:  {3: 2, 4: 4, 5: 6, 6: 15},
	7:  {3: 2, 4: 4, 5: 6, 6: 15},
	8:  {3: 1, 4: 2, 5: 4, 6: 10},
	9:  {3: 1, 4: 2, 5: 4, 6: 10},
	10: {3: 1, 4: 2, 5: 4, 6: 10},
}

var spinMultipliers = []int{1, 5, 10}

var spinCounts = map[int][]int{
	4: {15, 10, 5},
	5: {19, 14, 9},
	6: {23, 18, 13},
}

var spinCountsByScatters = map[int][]models.BonusChoice{
	4: {
		{Spins: 15, Multiplier: 1},
		{Spins: 10, Multiplier: 5},
		{Spins: 5, Multiplier: 10},
	},
	5: {
		{Spins: 19, Multiplier: 1},
		{Spins: 14, Multiplier: 5},
		{Spins: 9, Multiplier: 10},
	},
	6: {
		{Spins: 23, Multiplier: 1},
		{Spins: 18, Multiplier: 5},
		{Spins: 13, Multiplier: 10},
	},
}

var mapSpinsScatters = map[int]int{
	15: 4,
	10: 4,
	5:  4,
	19: 5,
	14: 5,
	9:  5,
	23: 6,
	18: 6,
	13: 6,
}

// scatter multiplier map
var smp = map[int]float64{
	4: 1,
	5: 1.67094,
	6: 2.46676,
}

var additionalSpinsTriggerStops = []int{0, 0, 0, 0, 0, 0, 0}

// need to be sure that we have enough space between scatter symbols
func findScatterStops(reels [][]int, topReel []int) []int {
	result := make([]int, len(reels))

	added := 0
	for i, reel := range reels {
		_, index, ok := lo.FindIndexOf(reel, func(item int) bool {
			return item == scatterSymbol
		})

		if added >= buyBonusScatterQty {
			newIndex := index - 1
			if newIndex < 0 {
				newIndex = len(reel) - 1
			}

			result[i] = newIndex
			continue
		}

		if ok {
			result[i] = index + rand.Intn(2)
			added++
		}
	}

	result = append(result, rand.Intn(len(topReel)))

	return result
}
