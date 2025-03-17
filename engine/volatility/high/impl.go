package high

import (
	vol "aztec-pyramids/engine/volatility"

	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine/utils/volatility"
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
)

var Volatility high

type high struct {
	volatility.High
}

func (m high) Config(rand rng.Client, rtp float64) *vol.Config {
	return vol.NewConfig(
		rand, rtp, mainRTPPercentage,
		availableReels, availableTopReels,
		lowRTPBase, highRTPBase, lowRTPAnte, highRTPAnte, lowBuyBonus, highBuyBonus,
		lowBuyBonusCorrectionCoef, highBuyBonusCorrectionCoef,
		configBaseLow, configBaseHigh, configBonusLow, configBonusHigh, configAnteLow, configAnteHigh,
		rtpTable, randomRtpTable,
		baseLowBonusRate, baseHighBonusRate, anteLowBonusRate, anteHighBonusRate,
		baseLow4ScatterCount, baseLow5ScatterCount, baseLow6ScatterCount,
		baseHigh4ScatterCount, baseHigh5ScatterCount, baseHigh6ScatterCount,
		anteLow4ScatterCount, anteLow5ScatterCount, anteLow6ScatterCount,
		anteHigh4ScatterCount, anteHigh5ScatterCount, anteHigh6ScatterCount,
		buyBonusStops,
	)
}
