package vol

import (
	"aztec-pyramids/engine/models"

	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"bitbucket.org/play-workspace/base-slot-server/utils"
	"github.com/samber/lo"
)

type Config struct {
	MainRTPPercentage float64
	AnteRTPPercentage float64

	AvailableReels        [][][]int
	AvailableTopReels     [][]int
	AvailableIndexedReels [][]map[int]int

	LowRTPBase, HighRTPBase   float64
	LowRTPAnte, HighRTPAnte   float64
	LowBuyBonus, HighBuyBonus float64

	LowBuyBonusCorrectionCoef  float64
	HighBuyBonusCorrectionCoef float64

	BuyBonusStops [][]int

	BaseConfigLowProb float64
	AnteConfigLowProb float64

	CfgBaseLow  *utils.Chooser[int, int]
	CfgBaseHigh *utils.Chooser[int, int]

	CfgBonusLow  *utils.Chooser[int, int]
	CfgBonusHigh *utils.Chooser[int, int]

	CfgAnteLow  *utils.Chooser[int, int]
	CfgAnteHigh *utils.Chooser[int, int]

	RtpTable       map[models.BonusChoice]models.BonusRtp
	RandomRtpTable map[int]models.BonusRtp

	BaseLowBonusRate  float64
	BaseHighBonusRate float64

	AnteLowBonusRate  float64
	AnteHighBonusRate float64

	BaseLow4ScatterCount  float64
	BaseLow5ScatterCount  float64
	BaseLow6ScatterCount  float64
	BaseHigh4ScatterCount float64
	BaseHigh5ScatterCount float64
	BaseHigh6ScatterCount float64

	AnteLow4ScatterCount  float64
	AnteLow5ScatterCount  float64
	AnteLow6ScatterCount  float64
	AnteHigh4ScatterCount float64
	AnteHigh5ScatterCount float64
	AnteHigh6ScatterCount float64
}

func NewConfig(rand rng.Client, rtp, mainRTPPercentage float64,
	availableReels [][][]int, availableTopReels [][]int,
	lowRTPBase, highRTPBase, lowRTPAnte, highRTPAnte, lowBuyBonus, highBuyBonus float64,
	lowBuyBonusCorrectionCoef, highBuyBonusCorrectionCoef float64,
	configBaseLow, configBaseHigh, configBonusLow, configBonusHigh, configAnteLow, configAnteHigh map[int]int,
	rtpTable map[models.BonusChoice]models.BonusRtp,
	randomRtpTable map[int]models.BonusRtp,
	baseLowBonusRate, baseHighBonusRate, anteLowBonusRate, anteHighBonusRate,
	baseLow4ScatterCount, baseLow5ScatterCount, baseLow6ScatterCount,
	baseHigh4ScatterCount, baseHigh5ScatterCount, baseHigh6ScatterCount,
	anteLow4ScatterCount, anteLow5ScatterCount, anteLow6ScatterCount,
	anteHigh4ScatterCount, anteHigh5ScatterCount, anteHigh6ScatterCount float64,
	buyBonusStops [][]int,
) *Config {
	if len(availableReels) != len(availableTopReels) {
		panic("len(availableReels) != len(availableTopReels)")
	}

	cfgBaseLow, err := utils.NewChooserFromMap(rand, configBaseLow)
	if err != nil {
		panic(err)
	}

	cfgBaseHigh, err := utils.NewChooserFromMap(rand, configBaseHigh)
	if err != nil {
		panic(err)
	}

	cfgBonusLow, err := utils.NewChooserFromMap(rand, configBonusLow)
	if err != nil {
		panic(err)
	}

	cfgBonusHigh, err := utils.NewChooserFromMap(rand, configBonusHigh)
	if err != nil {
		panic(err)
	}

	cfgAnteLow, err := utils.NewChooserFromMap(rand, configAnteLow)
	if err != nil {
		panic(err)
	}

	cfgAnteHigh, err := utils.NewChooserFromMap(rand, configAnteHigh)
	if err != nil {
		panic(err)
	}

	cfg := &Config{
		MainRTPPercentage: mainRTPPercentage,
		AnteRTPPercentage: mainRTPPercentage / 1.25,

		AvailableReels:    availableReels,
		AvailableTopReels: availableTopReels,
		AvailableIndexedReels: lo.Map(availableReels, func(reel [][]int, i int) []map[int]int {
			return createIndexedReels(reel, availableTopReels[i])
		}),

		LowRTPBase:   lowRTPBase,
		HighRTPBase:  highRTPBase,
		LowRTPAnte:   lowRTPAnte,
		HighRTPAnte:  highRTPAnte,
		LowBuyBonus:  lowBuyBonus,
		HighBuyBonus: highBuyBonus,

		LowBuyBonusCorrectionCoef:  lowBuyBonusCorrectionCoef,
		HighBuyBonusCorrectionCoef: highBuyBonusCorrectionCoef,

		CfgBaseLow:   cfgBaseLow,
		CfgBaseHigh:  cfgBaseHigh,
		CfgBonusLow:  cfgBonusLow,
		CfgBonusHigh: cfgBonusHigh,
		CfgAnteLow:   cfgAnteLow,
		CfgAnteHigh:  cfgAnteHigh,

		RtpTable:       rtpTable,
		RandomRtpTable: randomRtpTable,

		BaseLowBonusRate:  baseLowBonusRate,
		BaseHighBonusRate: baseHighBonusRate,

		AnteLowBonusRate:  anteLowBonusRate,
		AnteHighBonusRate: anteHighBonusRate,

		BaseLow4ScatterCount:  baseLow4ScatterCount,
		BaseLow5ScatterCount:  baseLow5ScatterCount,
		BaseLow6ScatterCount:  baseLow6ScatterCount,
		BaseHigh4ScatterCount: baseHigh4ScatterCount,
		BaseHigh5ScatterCount: baseHigh5ScatterCount,
		BaseHigh6ScatterCount: baseHigh6ScatterCount,

		AnteLow4ScatterCount:  anteLow4ScatterCount,
		AnteLow5ScatterCount:  anteLow5ScatterCount,
		AnteLow6ScatterCount:  anteLow6ScatterCount,
		AnteHigh4ScatterCount: anteHigh4ScatterCount,
		AnteHigh5ScatterCount: anteHigh5ScatterCount,
		AnteHigh6ScatterCount: anteHigh6ScatterCount,

		BuyBonusStops: buyBonusStops,
	}

	cfg.setRTP(rtp)

	return cfg
}

const (
	lenWithTopReel = 7
	topReelIndex   = 6
)

func createIndexedReels(reelSet [][]int, topReelSet []int) []map[int]int {
	result := make([]map[int]int, lenWithTopReel)

	for i, reel := range reelSet {
		result[i] = make(map[int]int, len(reel))
		for symbolIndex, symbol := range reel {
			result[i][symbolIndex] = symbol
		}
	}

	result[topReelIndex] = make(map[int]int, len(topReelSet))
	for symbolIndex, symbol := range topReelSet {
		result[topReelIndex][symbolIndex] = symbol
	}

	return result
}

func (c *Config) setRTP(rtp float64) {
	c.BaseConfigLowProb = c.calcBaseRTP(rtp)
	c.AnteConfigLowProb = c.calcAnteRTP(rtp)
}

func (c *Config) CalcRTP(rtp float64) (float64, float64) {
	return c.calcBaseRTP(rtp), c.calcAnteRTP(rtp)
}

func (c *Config) calcBaseRTP(rtp float64) float64 {
	mainRTP := rtp * c.MainRTPPercentage

	if mainRTP > c.HighRTPBase {
		return 0
	}

	if mainRTP < c.LowRTPBase {
		return 1
	}

	return (c.HighRTPBase - mainRTP) / (c.HighRTPBase - c.LowRTPBase)
}

func (c *Config) calcAnteRTP(rtp float64) float64 {
	anteRTP := rtp * c.AnteRTPPercentage

	if anteRTP > c.HighRTPAnte {
		return 0
	}

	if anteRTP < c.LowRTPAnte {
		return 1
	}

	return (c.HighRTPAnte - anteRTP) / (c.HighRTPAnte - c.LowRTPAnte)
}

func (c *Config) GetBuyBonusCorrectionCoef(value, rtp float64) float64 {
	if rtp > c.HighBuyBonus {
		return c.HighBuyBonusCorrectionCoef
	}

	if rtp < c.LowBuyBonus {
		return c.LowBuyBonusCorrectionCoef
	}

	if value > (c.HighBuyBonus-rtp)/(c.HighBuyBonus-c.LowBuyBonus) {
		return c.HighBuyBonusCorrectionCoef
	}

	return c.LowBuyBonusCorrectionCoef
}
