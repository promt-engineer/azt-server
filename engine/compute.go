package engine

import (
	"aztec-pyramids/engine/models"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine/utils"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine/utils/volatility"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/errs"
	baseUtils "bitbucket.org/play-workspace/base-slot-server/utils"
	"fmt"
	"github.com/samber/lo"
)

func (s *SpinFactory) computeBasicWindow(
	ag AwardGetter, reels []map[int]int, stops []int, config *baseUtils.Chooser[int, int],
) (award int64, avalanches []utils.Avalanche[int], bonusChoice []models.BonusChoice, err error) {
	window := Window{config: config}

	window.changeWindowSize = true

	avalanches, award, err = utils.Spin[int](ag, &window, reels, stops)
	if err != nil {
		return 0, nil, nil, err
	}

	if window.GetScatterQty() >= needScattersToTriggerBonus {
		avalanches[len(avalanches)-1].PayItems = append(avalanches[len(avalanches)-1].PayItems, utils.PayItem[int]{
			Symbol:  scatterSymbol,
			Indexes: window.GetIndexesBySymbol(scatterSymbol),
		})

		WinsMap.Inc(window.GetScatterQty())

		bonusChoice, err := s.generateBonusChoice(window.GetScatterQty())
		if err != nil {
			return 0, nil, nil, err
		}

		return award, avalanches, bonusChoice, nil
	}

	return award, avalanches, nil, nil
}

func (s *SpinFactory) computeBonusWindow(
	ag AwardGetter, bonusChoice models.BonusChoice, cheats *Cheats,
	isAnteBet, isBuyBonus bool, rtp float64, vol *volatility.Type,
) (*Bonus, error) {
	config, err := s.getBonusConfig(rtp, isAnteBet, isBuyBonus, bonusChoice, vol)
	if err != nil {
		return nil, err
	}

	var (
		bonus      = new(Bonus)
		window     = Window{config: config}
		extraSpins = 0
	)

	award := int64(0)

	for count := bonusChoice.Spins; count > 0 || extraSpins > 0; count-- {
		var multipliers []int

		stops, err := s.getStops(cheats, reelSet3Code, vol)
		if err != nil {
			return nil, err
		}

		window.changeWindowSize = true

		indexedReelSet3 := s.Cfg(vol).AvailableIndexedReels[reelSet3Code]

		avalanches, _, err := utils.Spin[int](ag, &window, indexedReelSet3, stops)
		if err != nil {
			return nil, err
		}

		spinAward := int64(0)

		mapAvalanches := lo.Map(avalanches, func(item utils.Avalanche[int], index int) models.Avalanche {
			payItems := make([]models.PayItem, len(item.PayItems))

			for i, payItem := range item.PayItems {
				awardWithMultiplier := payItem.Award * int64(bonusChoice.Multiplier)
				payItems[i] = models.PayItem{
					Symbol:              payItem.Symbol,
					Indexes:             payItem.Indexes,
					Award:               payItem.Award,
					Multiplier:          bonusChoice.Multiplier,
					AwardWithMultiplier: awardWithMultiplier,
				}
				spinAward += awardWithMultiplier
			}

			multipliers = append(multipliers, bonusChoice.Multiplier)

			if len(payItems) > 0 {
				bonusChoice.Multiplier++
			}

			return models.Avalanche{
				Window:   item.Window,
				PayItems: payItems,
			}
		})

		if extraSpins > 0 {
			count += extraSpins
			extraSpins = 0
		}

		if window.GetScatterQty() >= needScattersToTriggerExtraSpins {
			mapAvalanches[len(mapAvalanches)-1].PayItems = append(mapAvalanches[len(mapAvalanches)-1].PayItems, models.PayItem{
				Symbol:  scatterSymbol,
				Indexes: window.GetIndexesBySymbol(scatterSymbol),
			})

			extraSpins = extraSpinCount

			window.SetScatterQty(0)
		}

		bonus.Spins = append(bonus.Spins, SpinBonus{
			Stops:               stops,
			Award:               spinAward,
			BonusSpinsLeft:      count,
			BonusSpinsTriggered: extraSpins,
			Multiplier:          multipliers,
			Avalanches:          mapAvalanches,
		})

		award += spinAward
	}

	bonus.Win = award

	return bonus, nil
}

func (s *SpinFactory) generateBonusChoice(scatterQty int) ([]models.BonusChoice, error) {
	var bonusChoice = make([]models.BonusChoice, 4)

	copy(bonusChoice, spinCountsByScatters[scatterQty])

	randomSpins, err := s.rand.Rand(uint64(len(spinCounts[scatterQty])))
	if err != nil {
		return nil, err
	}

	randomMultipliers, err := s.rand.Rand(uint64(len(spinMultipliers)))
	if err != nil {
		return nil, err
	}

	randomBonusChoice := models.BonusChoice{
		Spins:      spinCounts[scatterQty][randomSpins],
		Multiplier: spinMultipliers[randomMultipliers],
		Random:     true,
	}

	bonusChoice[len(bonusChoice)-1] = randomBonusChoice

	return bonusChoice, nil
}

func (s *SpinFactory) getStops(cheats *Cheats, reelCode int, vol *volatility.Type) ([]int, error) {
	reel := s.Cfg(vol).AvailableReels[reelCode]

	if cheats != nil {
		if reelCode == reelSet3Code {
			if cheats.AdditionalTriggerCount > 0 {
				cheats.AdditionalTriggerCount--

				stops := make([]int, len(additionalSpinsTriggerStops))
				copy(stops, additionalSpinsTriggerStops)

				return stops, nil
			}
		} else {
			if len(cheats.Stops) != len(reel)+1 {
				return nil, errs.ErrBadDataGiven
			}

			last := len(cheats.Stops) - 1

			for i := range cheats.Stops {
				if i != last && cheats.Stops[i] > len(reel[i]) {
					return nil, errs.ErrBadDataGiven
				}

				if cheats.Stops[last] > len(s.Cfg(vol).AvailableTopReels[reelCode]) {
					return nil, errs.ErrBadDataGiven
				}
			}

			return cheats.Stops, nil
		}
	}

	req := lo.Map(reel, func(item []int, index int) uint64 {
		return uint64(len(item))
	})

	req = append(req, uint64(len(s.Cfg(vol).AvailableTopReels[reelCode])))

	res, err := s.rand.RandSlice(req)
	if err != nil {
		return nil, err
	}

	return lo.Map(res, func(item uint64, index int) int {
		return int(item)
	}), nil
}

func (s *SpinFactory) getBaseConfig(anteBet bool, rtp *int64, vol *volatility.Type) (*baseUtils.Chooser[int, int], error) {
	value, err := s.rand.RandFloat()
	if err != nil {
		return nil, err
	}

	if rtp != nil {
		userBaseConfigLowProb, userAnteConfigLowProb := s.Cfg(vol).CalcRTP(float64(*rtp))

		return s.chooseConfig(anteBet, value, userBaseConfigLowProb, userAnteConfigLowProb, vol), nil
	}

	return s.chooseConfig(anteBet, value, s.Cfg(vol).BaseConfigLowProb, s.Cfg(vol).AnteConfigLowProb, vol), nil
}

func (s *SpinFactory) getBonusConfig(rtp float64, isAnteBet, isBuyBonus bool, bc models.BonusChoice, vol *volatility.Type) (*baseUtils.Chooser[int, int], error) {
	avgFreeWin := s.calcAvgFreeWin(rtp, isAnteBet, isBuyBonus, vol)
	correctionCof, err := s.calcCorrectionCof(isAnteBet, isBuyBonus, vol, rtp)
	if err != nil {
		return nil, err
	}

	scatterCount := mapSpinsScatters[bc.Spins]

	var realAvgFreeWin float64

	mlt, ok := smp[scatterCount]
	if !ok {
		return nil, fmt.Errorf("invalid scatter count: %d", scatterCount)
	}

	realAvgFreeWin = mlt * avgFreeWin / correctionCof

	var rtpData models.BonusRtp

	if bc.Random {
		rtpData = s.Cfg(vol).RandomRtpTable[scatterCount]
	} else {
		rtpData = s.Cfg(vol).RtpTable[bc]
	}

	configFreeLowProb := rtpData.CalcLowProb(realAvgFreeWin)

	value, err := s.rand.RandFloat()
	if err != nil {
		return nil, err
	}

	if value < configFreeLowProb {
		return s.Cfg(vol).CfgBonusLow, nil
	}

	return s.Cfg(vol).CfgBonusHigh, nil
}

func (s *SpinFactory) chooseConfig(anteBet bool, value, baseLow, anteLow float64, vol *volatility.Type) *baseUtils.Chooser[int, int] {
	if anteBet {
		if value <= anteLow {
			return s.Cfg(vol).CfgAnteLow
		}

		return s.Cfg(vol).CfgAnteHigh
	}

	if value <= baseLow {
		return s.Cfg(vol).CfgBaseLow
	}

	return s.Cfg(vol).CfgBaseHigh
}

func (s *SpinFactory) calcAvgFreeWin(rtp float64, isAnteBet, isBuyBonus bool, vol *volatility.Type) float64 {
	var (
		freeRTP                     float64
		lowProb, highProb           float64
		lowBonusRate, highBonusRate float64
	)

	switch {
	case isAnteBet:
		freeRTP = s.calcFreeRTP(rtp, true, vol)
		lowProb = s.Cfg(vol).AnteConfigLowProb
		lowBonusRate = s.Cfg(vol).AnteLowBonusRate
		highBonusRate = s.Cfg(vol).AnteHighBonusRate
	case isBuyBonus:
		return avgFreeWinBuyBonus
	default:
		freeRTP = s.calcFreeRTP(rtp, false, vol)
		lowProb = s.Cfg(vol).BaseConfigLowProb
		lowBonusRate = s.Cfg(vol).BaseLowBonusRate
		highBonusRate = s.Cfg(vol).BaseHighBonusRate
	}

	highProb = 1 - lowProb

	return freeRTP / (lowProb/lowBonusRate + highProb/highBonusRate)
}

func (s *SpinFactory) calcFreeRTP(rtp float64, isAnteBet bool, vol *volatility.Type) float64 {
	if isAnteBet {
		return (rtp - rtp*s.Cfg(vol).AnteRTPPercentage) / 100
	}

	return (rtp - rtp*s.Cfg(vol).MainRTPPercentage) / 100
}

func (s *SpinFactory) calcCorrectionCof(isAnteBet, isBuyBonus bool, vol *volatility.Type, rtp float64) (float64, error) {
	var (
		lowProb, highProb                   float64
		low4ScatterCount, high4ScatterCount float64
		low5ScatterCount, high5ScatterCount float64
		low6ScatterCount, high6ScatterCount float64
	)

	switch {
	case isAnteBet:
		lowProb = s.Cfg(vol).AnteConfigLowProb
		highProb = 1 - lowProb

		low4ScatterCount = s.Cfg(vol).AnteLow4ScatterCount
		low5ScatterCount = s.Cfg(vol).AnteLow5ScatterCount
		low6ScatterCount = s.Cfg(vol).AnteLow6ScatterCount

		high4ScatterCount = s.Cfg(vol).AnteHigh4ScatterCount
		high5ScatterCount = s.Cfg(vol).AnteHigh5ScatterCount
		high6ScatterCount = s.Cfg(vol).AnteHigh6ScatterCount
	case isBuyBonus:
		value, err := s.rand.RandFloat()
		if err != nil {
			return 0, err
		}

		return s.Cfg(vol).GetBuyBonusCorrectionCoef(value, rtp), nil
	default:
		lowProb = s.Cfg(vol).BaseConfigLowProb
		highProb = 1 - lowProb

		low4ScatterCount = s.Cfg(vol).BaseLow4ScatterCount
		low5ScatterCount = s.Cfg(vol).BaseLow5ScatterCount
		low6ScatterCount = s.Cfg(vol).BaseLow6ScatterCount

		high4ScatterCount = s.Cfg(vol).BaseHigh4ScatterCount
		high5ScatterCount = s.Cfg(vol).BaseHigh5ScatterCount
		high6ScatterCount = s.Cfg(vol).BaseHigh6ScatterCount
	}

	mul4 := smp[4]
	mul5 := smp[5]
	mul6 := smp[6]

	return (lowProb*(low4ScatterCount*mul4+low5ScatterCount*mul5+low6ScatterCount*mul6) + highProb*(high4ScatterCount*mul4+high5ScatterCount*mul5+high6ScatterCount*mul6)) /
		(lowProb*(low4ScatterCount+low5ScatterCount+low6ScatterCount) + highProb*(high4ScatterCount+high5ScatterCount+high6ScatterCount)), nil
}
