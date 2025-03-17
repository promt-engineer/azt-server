package engine

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/engine"
	"encoding/json"
	"go.uber.org/zap"
)

type RestoringIndexes struct {
	BaseSpinIndex   int  `json:"base_spin_index"`
	BonusSpinIndex  int  `json:"bonus_spin_index"`
	GambleCollected bool `json:"gamble_collected"`
}

func (r *RestoringIndexes) IsShown(spin engine.Spin) bool {
	spinTyped, ok := spin.(*SpinBase)

	if !ok {
		zap.S().Error("can not parse spin")

		return true
	}

	if !spinTyped.BonusTriggered() {
		if len(spinTyped.BonusChoice) == 0 {
			return r.BaseSpinIndex == 1
		}

		// bonus is not triggered, but bonus choice is not empty
		return false
	}

	return r.BaseSpinIndex == 1 && r.BonusSpinIndex == len(spinTyped.Bonus.Spins)
}

func (r *RestoringIndexes) Update(payload interface{}) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, r)
}
