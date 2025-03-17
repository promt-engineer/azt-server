package engine

import (
	"aztec-pyramids/engine/models"
	"go.uber.org/zap"

	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine/utils"
)

type SpinBase struct {
	WagerVal        int64  `json:"wager"`
	AnteBetWagerVal *int64 `json:"ante_bet_wager"`
	Win             int64  `json:"award"`
	Stops           []int  `json:"stops"`
	BuyBonus        bool   `json:"buy_bonus"`
	AnteBet         bool   `json:"ante_bet"`

	Avalanches []utils.Avalanche[int] `json:"avalanches"`

	BonusChoice []models.BonusChoice `json:"bonus_choice"`

	Bonus *Bonus `json:"bonus"`

	Gamble engine.Gamble `json:"gambles,omitempty"`
}

func (s *SpinBase) BaseAward() int64 {
	return s.Win
}

func (s *SpinBase) BonusAward() int64 {
	if s.Bonus == nil {
		return 0
	}

	return s.Bonus.Award()
}

func (s *SpinBase) GambleAward() int64 {
	return 0
}

func (s *SpinBase) OriginalWager() int64 {
	return s.WagerVal
}

func (s *SpinBase) Wager() int64 {
	if s.AnteBet && s.AnteBetWagerVal != nil {
		return *s.AnteBetWagerVal
	}

	return s.WagerVal
}

func (s *SpinBase) BonusTriggered() bool {
	return s.Bonus != nil
}

func (s *SpinBase) GambleQuantity() int {
	return 0
}

func (s *SpinBase) GetGamble() *engine.Gamble {
	return &s.Gamble
}

func (s *SpinBase) CanGamble(restoringIndexes engine.RestoringIndexes) bool {
	if s.Bonus != nil || len(s.BonusChoice) > 0 {
		return false
	}

	restoringIndexesTyped, ok := restoringIndexes.(*RestoringIndexes)
	if !ok {
		zap.S().Error("can not parse restoring indexes")

		return true
	}

	return !restoringIndexesTyped.GambleCollected
}

func (s *SpinBase) DeepCopy() engine.Spin {
	spin := &SpinBase{
		WagerVal:        s.WagerVal,
		AnteBetWagerVal: s.AnteBetWagerVal,
		Win:             s.Win,
		BuyBonus:        s.BuyBonus,
		AnteBet:         s.AnteBet,

		Stops:       make([]int, len(s.Stops)),
		Avalanches:  make([]utils.Avalanche[int], len(s.Avalanches)),
		BonusChoice: make([]models.BonusChoice, len(s.BonusChoice)),
	}

	copy(spin.Stops, s.Stops)
	copy(spin.Avalanches, s.Avalanches)
	copy(spin.BonusChoice, s.BonusChoice)

	if s.Bonus != nil {
		spin.Bonus = s.Bonus.deepCopy()
	}

	return spin
}
