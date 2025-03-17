package engine

import (
	"aztec-pyramids/engine/models"
	vol "aztec-pyramids/engine/volatility"
	"aztec-pyramids/engine/volatility/high"
	"aztec-pyramids/engine/volatility/low"
	"aztec-pyramids/engine/volatility/medium"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine/utils/volatility"
	"bitbucket.org/play-workspace/base-slot-server/pkg/rng"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"math/rand"
)

func Bootstrap(rand rng.Client, volType volatility.Type, rtp float64) *engine.Bootstrap {
	cfgMap := volatility.NewVolatilityMap[vol.Config](rand, rtp, high.Volatility, medium.Volatility, low.Volatility)

	factory := NewSpinFactory(rand, cfgMap, volType, rtp)

	return &engine.Bootstrap{
		SpinFactory:   factory,
		HTTPTransport: true,

		FreeSpinsFeature:    true,
		GambleAnyWinFeature: true,

		AnteBetMultiplier: anteBetMultiplier,

		HistoryHandlingType: engine.SequentialRestoring,

		EngineInfo: map[string]interface{}{
			"buy_bonus_multiplier": buyBonusMultiplier,
		},
	}
}

type SpinFactory struct {
	rand rng.Client

	cfgMap            volatility.ConfigMap[vol.Config]
	defaultVolatility volatility.Type
	defaultRTP        float64
}

func NewSpinFactory(
	rand rng.Client, cfgMap volatility.ConfigMap[vol.Config], vol volatility.Type, rtp float64,
) *SpinFactory {
	return &SpinFactory{
		rand:              rand,
		cfgMap:            cfgMap,
		defaultVolatility: vol,
		defaultRTP:        rtp,
	}
}

func (s *SpinFactory) Cfg(vol *volatility.Type) *vol.Config {
	if vol != nil {
		return s.cfgMap[*vol]
	}

	return s.cfgMap[s.defaultVolatility]
}

func (s *SpinFactory) Generate(
	ctx engine.Context, wager int64, parameters interface{},
) (engine.Spin, engine.RestoringIndexes, error) {
	isAnteBet, isBuyBonus, err := getFeaturesFromParams(parameters)
	if err != nil {
		return nil, nil, err
	}

	rtp, volType, err := getUserParams(ctx)
	if err != nil {
		return nil, nil, err
	}
	//zap.S().Info("Generating...")
	//if rtp != nil {
	//	zap.S().Info("RTP: ", *rtp)
	//}
	//if volType != nil {
	//	zap.S().Info("Volatility: ", *volType)
	//}

	cheats, err := getCheatsFromCtx(ctx)
	if err != nil {
		return nil, nil, err
	}

	ag, err := NewAwardGetter(wager)
	if err != nil {
		return nil, nil, err
	}

	var anteBetWagerVal *int64

	reelCode := reelSet1Code
	if isAnteBet {
		reelCode = reelSet2Code

		anteBetWagerVal = lo.ToPtr(wager * anteBetMultiplier / anteBetDivider)
	}

	var stops []int
	if isBuyBonus {
		reelCode = reelSet3Code

		buyBonusStops := s.Cfg(volType).BuyBonusStops

		stops = buyBonusStops[rand.Intn(len(buyBonusStops))]

		wager *= buyBonusMultiplier
	} else {
		stops, err = s.getStops(cheats, reelCode, volType)
		if err != nil {
			return nil, nil, err
		}
	}

	config, err := s.getBaseConfig(isAnteBet, rtp, volType)
	if err != nil {
		return nil, nil, err
	}

	award, avalanches, bonusChoice, err := s.computeBasicWindow(ag, s.Cfg(volType).AvailableIndexedReels[reelCode], stops, config)
	if err != nil {
		return nil, nil, err
	}

	return &SpinBase{
		WagerVal:        wager,
		AnteBetWagerVal: anteBetWagerVal,
		Win:             award,
		Stops:           stops,
		AnteBet:         isAnteBet,
		BuyBonus:        isBuyBonus,

		Avalanches: avalanches,

		BonusChoice: bonusChoice,
	}, &RestoringIndexes{}, nil
}

func (s *SpinFactory) KeepGenerate(ctx engine.Context, parameters interface{}) (engine.Spin, bool, error) {
	if ctx.LastSpin == nil {
		return nil, false, fmt.Errorf("last spin is nil")
	}

	if ctx.LastSpin.BonusTriggered() {
		return nil, false, fmt.Errorf("bonus already triggered")
	}

	spin, ok := ctx.LastSpin.(*SpinBase)
	if !ok {
		return nil, false, fmt.Errorf("can not parse spin")
	}

	if spin.BonusChoice == nil {
		return nil, false, fmt.Errorf("bonus is not available for this spin")
	}

	wager := spin.Wager()
	if spin.BuyBonus {
		wager /= buyBonusMultiplier
	}

	cheats, err := getCheatsFromCtx(ctx)
	if err != nil {
		return nil, false, err
	}

	ag, err := NewAwardGetter(wager)
	if err != nil {
		return nil, false, err
	}

	bc, err := UnmarshalTo[models.BonusChoice](parameters)
	if err != nil {
		return nil, false, err
	}

	if _, exist := lo.Find(spin.BonusChoice, func(item models.BonusChoice) bool {
		return item.Spins == bc.Spins && item.Multiplier == bc.Multiplier
	}); !exist {
		return nil, false, fmt.Errorf("invalid bonus choice")
	}

	rtp := s.defaultRTP

	rtpPtr, volType, err := getUserParams(ctx)
	if err != nil {
		return nil, false, err
	}

	if rtpPtr != nil {
		rtp = float64(*rtpPtr)
	}

	bonus, err := s.computeBonusWindow(ag, bc, cheats, spin.AnteBet, spin.BuyBonus, rtp, volType)
	if err != nil {
		return nil, false, err
	}

	spin.Bonus = bonus

	return spin, true, nil
}

func (s *SpinFactory) UnmarshalJSONSpin(bytes []byte) (engine.Spin, error) {
	spin := SpinBase{}
	err := json.Unmarshal(bytes, &spin)

	return &spin, err
}

func (s *SpinFactory) UnmarshalJSONRestoringIndexes(bytes []byte) (engine.RestoringIndexes, error) {
	restoringIndexes := RestoringIndexes{}
	err := json.Unmarshal(bytes, &restoringIndexes)

	return &restoringIndexes, err
}

func (s *SpinFactory) GetRngClient() rng.Client {
	return s.rand
}

func getUserParams(ctx engine.Context) (*int64, *volatility.Type, error) {
	if ctx.UserParams != nil {
		if ctx.UserParams.Volatility != nil {
			v, err := volatility.VolFromStr(*ctx.UserParams.Volatility)
			if err != nil {
				return nil, nil, err
			}

			return ctx.UserParams.RTP, &v, nil
		}

		return ctx.UserParams.RTP, nil, nil
	}

	return nil, nil, nil
}
